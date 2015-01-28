#!/bin/bash

which entr > /dev/null
if [ $? -ne 0 ]; then
    echo "installing entr"
    brew install entr
fi

go test ./...
while sleep 1; do
  ls *.go | entr -d -c go test ./...
done
