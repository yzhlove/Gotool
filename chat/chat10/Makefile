SHELL=/bin/bash

App:=docker-test
Env:=/var/run/docker.sock
Dockerfile:=Dockerfile
Arch:=arm64

start:
	@ echo -e "stop docker-test if running ...";\
	docker stop $(App) 2>/dev/null || true;\
	docker rm $(App) 2>/dev/null || true;\
	echo -e "start docker-test ...";\
	docker run --rm=true --name $(App) -h $(App) -it -p 8081:8081 -v $(Env):/app/docker.sock $(App)

build:
	@ echo -e "\nStart compile and rebuild...";\
	CGO_ENABLED=0 GOOS=linux GOARCH=$(Arch) go build -o $(App) .; \
	docker rm -f $(App) 2>/dev/null || true;\
	docker rmi $(App) 2>/dev/null || true;\
	docker build --platform linux/$(Arch) --no-cache -t $(App) -f $(Dockerfile) .;\
	rm $(App)

.PHONY: build start