#!/bin/bash

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