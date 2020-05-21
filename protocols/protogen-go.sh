#!/bin/bash
protoc -I/usr/local/include/ -I./user/ user/user.proto --go_out=plugins=grpc:../backend/pkg/api/
protoc -I/usr/local/include/ -I./hrtf/ hrtf/hrtf.proto --go_out=plugins=grpc:../backend/pkg/api/
