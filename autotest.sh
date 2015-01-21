#!/bin/bash
while sleep 1; do
  ls *.go | entr -c go test ./... -logtostderr=true
done
