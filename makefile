GATEWAY_SRC_PATH=internal/gateway
USER_SRC_PATH=internal/user
VIDEO_SRC_PATH=internal/video
FAVORITE_SRC_PATH=internal/favorite
FEED_SRC_PATH=internal/feed
COMMENT_SRC_PATH=internal/comment
RELATION_SRC_PATH=internal/relation
MESSAGE_SRC_PATH=internal/message
KITEX_GEN_PATH=sh/kitex-gen.sh

all: buildKitex buildUser buildVideo buildFavorite buildFeed buildComment buildRelation buildMessage buildGateway cleanData

buildKitex:
	sh $(KITEX_GEN_PATH)

buildUser:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(USER_SRC_PATH)/main $(USER_SRC_PATH)/*.go

buildVideo:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(VIDEO_SRC_PATH)/main $(VIDEO_SRC_PATH)/*.go

buildFavorite:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(FAVORITE_SRC_PATH)/main $(FAVORITE_SRC_PATH)/*.go

buildFeed:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(FEED_SRC_PATH)/main $(FEED_SRC_PATH)/*.go

buildComment:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(COMMENT_SRC_PATH)/main $(COMMENT_SRC_PATH)/*.go

buildRelation:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(RELATION_SRC_PATH)/main $(RELATION_SRC_PATH)/*.go

buildMessage:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(MESSAGE_SRC_PATH)/main $(MESSAGE_SRC_PATH)/*.go

buildGateway:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GATEWAY_SRC_PATH)/main $(GATEWAY_SRC_PATH)/*.go

cleanData:
	rm -rf ./data