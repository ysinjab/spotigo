#!/bin/sh
protoc -I=pkg/albums --go_out=pkg/albums pkg/albums/albums.proto
