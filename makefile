GATEWAY_SRC_PATH=internal/gateway
USER_SRC_PATH=internal/user
VIDEO_SRC_PATH=internal/video
FAVORITE_SRC_PATH=internal/favorite
COMMENT_SRC_PATH=internal/comment
RELATION_SRC_PATH=internal/relation
MESSAGE_SRC_PATH=internal/message
KITEX_GEN_PATH=sh/kitex-gen.sh
FFMPEG_SRC_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
FFMPEG_SRC_PATH=internal/video/ffmpeg
exist = $(shell if [ -f $(FFMPEG_SRC_PATH) ]; then echo "true"; else echo "false"; fi)

all: buildKitex buildUser buildVideo buildFavorite buildComment buildRelation buildMessage buildGateway download_ffmpeg

buildKitex:
	sh $(KITEX_GEN_PATH)

download_ffmpeg:
ifeq ("$(exist)", "false")
	@echo ">> downloading ffmpeg..."
	wget $(FFMPEG_SRC_URL) -O $(VIDEO_SRC_PATH)/ffmpeg.tar.xz
	tar -xvf $(VIDEO_SRC_PATH)/ffmpeg.tar.xz -C $(VIDEO_SRC_PATH)
	mv $(VIDEO_SRC_PATH)/ffmpeg-*-amd64-static/ffmpeg $(VIDEO_SRC_PATH)/ffmpeg
	rm -rf $(VIDEO_SRC_PATH)/ffmpeg-*-amd64-static
	rm -rf $(VIDEO_SRC_PATH)/ffmpeg.tar.xz
endif

buildUser:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(USER_SRC_PATH)/main $(USER_SRC_PATH)/*.go

buildVideo:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(VIDEO_SRC_PATH)/main $(VIDEO_SRC_PATH)/*.go

buildFavorite:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(FAVORITE_SRC_PATH)/main $(FAVORITE_SRC_PATH)/*.go

buildComment:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(COMMENT_SRC_PATH)/main $(COMMENT_SRC_PATH)/*.go

buildRelation:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(RELATION_SRC_PATH)/main $(RELATION_SRC_PATH)/*.go

buildMessage:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(MESSAGE_SRC_PATH)/main $(MESSAGE_SRC_PATH)/*.go

buildGateway:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GATEWAY_SRC_PATH)/main $(GATEWAY_SRC_PATH)/*.go