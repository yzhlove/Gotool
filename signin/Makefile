SHELL=/bin/bash

pb:
	@echo "生成ProtoBuf文件...";
	@protoc --proto_path=protocol/source \
		--go_opt=paths=source_relative \
		--go_out=protocol/proto protocol/source/*.proto \


.PHONY: pb