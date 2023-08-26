package main

import (
	"context"

	"toktik/internal/message/kitex_gen/message"
	"toktik/internal/message/pkg/ctx"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct {
	svcCtx *ctx.ServiceContext
}

func NewMessageServiceImpl(svcCtx *ctx.ServiceContext) *MessageServiceImpl {
	return &MessageServiceImpl{
		svcCtx: svcCtx,
	}
}

// ListMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) ListMessage(ctx context.Context, req *message.ListMessageReq) (resp *message.ListMessageRes, err error) {
	// TODO: Your code here...
	return
}

// SendMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) SendMessage(ctx context.Context, req *message.SendMessageReq) (resp *message.SendMessageRes, err error) {
	// TODO: Your code here...
	return
}

// GetLastMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) GetLastMessage(ctx context.Context, req *message.GetLastMessageReq) (resp *message.GetLastMessageRes, err error) {
	// TODO: Your code here...
	return
}
