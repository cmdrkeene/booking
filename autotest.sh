#!/bin/bash

which entr > /dev/null
if [ $? -ne 0 ]; then
    echo "installing entr"
    brew install entr
else
	echo "watching go files for testing"
fi

while sleep 1; do
  ls *.go | entr -d -c go test ./... -logtostderr=true
done
