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
-e CLIENT_IP=127.0.0.1 \
-p 6381:6381 \
-p 6382:6382 \
-p 6383:6383 \
-p 6384:6384 \
-p 6385:6385 \
-p 6386:6386 \
--net host \

build:
	@ echo -e "\nStart compile and rebuild...";\
 	echo -e "Arch is $(Arch)";\
	CGO_ENABLED=0 GOOS=linux GOARCH=$(Arch) go build -o $(App) .; \
	docker rm -f $(App);\
	docker rmi $(App);\
	docker build --platform linux/$(Arch) --no-cache -t $(App) -f $(Dockerfile) .;\
	rm $(App)

desc:
	@ echo "docker网络模式默认为桥接，无法找到宿主机的IP,如果需要 redis cluster-announce-ip 正确生效，需要指定--net host"; \
 	@ echo "使容器直接使用宿主机的网络栈，注意端口冲突!" \

start:
	@ echo -e "\nStart $(App)$(Group) with cmd...\n\ndocker run $(Cmd) $(App)\n"; \
	  docker run $(Cmd) $(App)

run: build start

.PHONY: build start run
