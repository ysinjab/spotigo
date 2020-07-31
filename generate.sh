#!/bin/sh
protoc -I pkg/albums albums.proto --go_out=plugins=grpc:pkg/albums