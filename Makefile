DOCKER_IMAGE = <my docker image>

export CGO_ENABLED = 0
export GOARCH = amd64
export GOOS = linux
export GO111MODULE = on
export GOFLAGS = -mod=vendor
export GOBIN = $(CURDIR)/bin



generate:
	protoc -I proto  --go_out=chat --go_opt=paths=source_relative \
	--go-grpc_out=chat --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false  chat.proto

build:
	go install -installsuffix "static" main/server.go

docker:
	docker build . -t $(DOCKER_IMAGE)

push:
	docker push $(DOCKER_IMAGE)

cert:
  mkdir -p cert
	openssl req -x509 -nodes -newkey rsa:2048 -days 365 -keyout cert/key.pem -out cert/cert.pem -subj "/CN=does.not.matter.com"

secret:
	kubectl create secret tls grpc-cert-self-managed --cert=cert/cert.pem --key=cert/key.pem

clean:
	$(RM) bin/server
	$(RM) chat/chat.pb.go
	$(RM) chat/chat_grpc.pb.go
	$(RM) cert/*
	docker rmi -f $(DOCKER_IMAGE)

.PHONY: cert