package main

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"sync"

	"gorm.io/gorm"

	"toktik/internal/favorite/kitex_gen/favorite"
	"toktik/internal/relation/kitex_gen/relation"
	"toktik/internal/user/kitex_gen/user"
	"toktik/internal/user/pkg/ctx"
	"toktik/internal/video/kitex_gen/video"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	svcCtx *ctx.ServiceContext
}

func NewUserServiceImpl(svcCtx *ctx.ServiceContext) *UserServiceImpl {
	return &UserServiceImpl{
		svcCtx: svcCtx,
	}
}

// Error messages.
var (
	errDuplicateUsername         = "用户名已被占用"
	errRegisterFailed            = "注册失败"
	errInvalidUsernameOrPassword = "用户名或密码错误"
	errLoginFailed               = "登录失败"
	errUserNotFound              = "用户不存在"
	errGetUserInfoFailed         = "获取用户信息失败"
)

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterRes, _ error) {
	resp = &user.RegisterRes{}

	err := checkRegisterParams(req.Username, req.Password)
	if err != nil {
		resp = &user.RegisterRes{
			Status: user.Status_ERROR,
			ErrMsg: err.Error(),
		}
		return
	}

	userId, err := s.svcCtx.UserService.CreateUser(req.Username, req.Password)
	if err != nil {
		resp.Status = user.Status_ERROR
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			resp.ErrMsg = errDuplicateUsername
		} else {
			resp.ErrMsg = errRegisterFailed
		}
		return
	}

	resp.UserId = userId

	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginRes, _ error) {
	resp = &user.LoginRes{}

	userId, err := s.svcCtx.UserService.Login(req.Username, req.Password)
	if err != nil {
		resp.Status = user.Status_ERROR
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.ErrMsg = errInvalidUsernameOrPassword
		} else {
			resp.ErrMsg = errLoginFailed
		}
		return
	}

	resp.UserId = userId

	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoRes, _ error) {
	resp = &user.GetUserInfoRes{}

	userInfo, err := s.svcCtx.UserService.GetUserById(req.ToUserId)
	if err != nil {
		resp.Status = user.Status_ERROR
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.ErrMsg = errUserNotFound
		} else {
			resp.ErrMsg = errGetUserInfoFailed
		}
		return
	}

	resp.User = &user.UserInfo{
		Id:              userInfo.Id,
		Name:            userInfo.Username,
		Avatar:          userInfo.Avatar,
		BackgroundImage: userInfo.BackgroundImage,
		Signature:       userInfo.Signature,
	}

	s.getInfoInOtherService(ctx, req.UserId, []*user.UserInfo{resp.User})

	return
}

// GetUserInfos implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfos(ctx context.Context, req *user.GetUserInfosReq) (resp *user.GetUserInfosRes, _ error) {
	resp = &user.GetUserInfosRes{}
	userInfos, err := s.svcCtx.UserService.GetUserByIds(req.ToUserIds)
	if err != nil {
		resp.Status = user.Status_ERROR
		resp.ErrMsg = err.Error()
		return
	}

	resp.Users = make([]*user.UserInfo, 0, len(userInfos))
	for _, userInfo := range userInfos {
		user := &user.UserInfo{
			Id:              userInfo.Id,
			Name:            userInfo.Username,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			Signature:       userInfo.Signature,
		}
		resp.Users = append(resp.Users, user)
	}

	s.getInfoInOtherService(ctx, req.UserId, resp.Users)

	return
}

func (s *UserServiceImpl) getInfoInOtherService(ctx context.Context, userId int64, userInfos []*user.UserInfo) {
	userId2UserInfo := make(map[int64]*user.UserInfo)
	userIdList := make([]int64, 0, len(userInfos))
	for _, userInfo := range userInfos {
		userId2UserInfo[userInfo.Id] = userInfo
		userIdList = append(userIdList, userInfo.Id)
	}
	wg := sync.WaitGroup{}

	// get user relation info
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.RelationClient.GetFollowInfo(ctx, &relation.GetFollowInfoReq{
			UserId:       userId,
			ToUserIdList: userIdList,
		})
		if err != nil || res.Status != relation.Status_OK {
			return
		}

		for _, followInfo := range res.FollowInfoList {
			userInfo := userId2UserInfo[followInfo.UserId]
			userInfo.FollowCount = followInfo.FollowCount
			userInfo.FollowerCount = followInfo.FollowerCount
			userInfo.IsFollow = followInfo.IsFollow
		}
	}()

	// get user favorite info
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.FavoriteClient.GetUserFavoriteInfo(ctx, &favorite.GetUserFavoriteInfoReq{
			UserIdList: userIdList,
		})
		if err != nil || res.Status != favorite.Status_OK {
			return
		}

		for _, favoriteInfo := range res.FavoriteInfoList {
			userInfo := userId2UserInfo[favoriteInfo.UserId]
			userInfo.FavoriteCount = favoriteInfo.FavoriteCount
			userInfo.TotalFavorited = favoriteInfo.TotalFavorited
		}
	}()

	// get user work count
	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := s.svcCtx.VideoClient.GetWorkCount(ctx, &video.GetWorkCountReq{
			UserIdList: userIdList,
		})
		if err != nil || res.Status != video.Status_OK {
			return
		}

		for _, workCountInfo := range res.WorkCountList {
			userInfo := userId2UserInfo[workCountInfo.UserId]
			userInfo.WorkCount = workCountInfo.WorkCount
		}
	}()

	wg.Wait()
}

var usernameRegexp *regexp.Regexp
var passwordRegexp *regexp.Regexp
var (
	errEmptyUsername        = errors.New("用户名不能为空")
	errEmptyPassword        = errors.New("密码不能为空")
	errShortPassword        = errors.New("密码长度不能小于6位")
	errLongUsername         = errors.New("用户名长度不能大于32位")
	errLongPassword         = errors.New("密码长度不能大于32位")
	errInvalidCharsUsername = errors.New("用户名只能包含中文、英文、数字、下划线")
	errInvalidCharsPassword = errors.New("密码只能包含英文、数字、特殊符号")
)

func checkRegisterParams(username string, password string) error {
	if username == "" {
		return errEmptyUsername
	}
	if password == "" {
		return errEmptyPassword
	}
	if len(password) < 6 {
		return errShortPassword
	}
	if len(username) > 32 {
		return errLongUsername
	}
	if len(password) > 32 {
		return errLongPassword
	}

	var err error
	if usernameRegexp == nil {
		usernameRegexp, err = regexp.Compile("^[\u4e00-\u9fa5_a-zA-Z0-9]+$")
		if err != nil {
			return err
		}
	}
	if !usernameRegexp.MatchString(username) {
		return errInvalidCharsUsername
	}

	if passwordRegexp == nil {
		passwordRegexp, err = regexp.Compile(`^[a-zA-Z0-9.~!@#$%^&*\-_+?]+$`)
		if err != nil {
			return err
		}
	}
	if !passwordRegexp.MatchString(password) {
		return errInvalidCharsPassword
	}

	return nil
}
