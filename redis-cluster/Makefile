SHELL=/bin/bash

App:=redis-cluster
Ports:="6381,6382,6383,6384,6385,6386"
Dockerfile:=Dockerfile
Arch:=$(shell uname -m)

ifneq ($(Arch),arm64)
	Arch=amd64
endif

Cmd :=--rm=true --name $(App) -h $(App) -it \
-e REDIS_PORTS=$(Ports) \
-p 6381:6381 \
-p 6382:6382 \
-p 6383:6383 \
-p 6384:6384 \
-p 6385:6385 \
-p 6386:6386

build:
	@ echo -e "\nStart compile and rebuild...";\
 	echo -e "Arch is $(Arch)";\
	CGO_ENABLED=0 GOOS=linux GOARCH=$(Arch) go build -o $(App) .; \
	docker rm -f $(App);\
	docker rmi $(App);\
	docker build --platform linux/$(Arch) --no-cache -t $(App) -f $(Dockerfile) .;\
	rm $(App)

start:
	@ echo -e "\nStart $(App)$(Group) with cmd...\n\ndocker run $(Cmd) $(App)\n"; \
	  docker run $(Cmd) $(App)

run: build start

.PHONY: build start run
