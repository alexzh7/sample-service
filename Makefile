.PHONY: proto

proto:
	protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative dvdstore.proto
