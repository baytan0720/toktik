#!/bin/bash

# env
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
kitex -version || exit

if [[ $(basename "$PWD") == "sh" ]]; then
    cd ..
elif [[ $(basename "$PWD") != "toktik" ]]; then
    echo "please run this script in toktik project root directory"
    exit 1
fi

mkdir -p "internal"
cd "internal" || exit

# comment
mkdir -p "comment"
cd "comment" || exit
kitex -I ../../idl -module toktik -service comment comment-service.proto
cd ..

# relation
mkdir -p "relation"
cd "relation" || exit
kitex -I ../../idl -module toktik -service relation relation-service.proto
cd ..

# relation
mkdir -p "user"
cd "user" || exit
kitex -I ../../idl -module toktik -service user user-service.proto
cd ..