#!/bin/bash

GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -ldflags="-w -s" -trimpath -o bin/prochide.dll
