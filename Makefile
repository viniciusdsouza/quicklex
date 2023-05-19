.PHONY: create-file build transfer-binary

default: build-cli

create-file:
	mkdir /opt/ql

build:
	go build -o ql .

transfer-binary:
	cp ./ql /opt/ql

build-cli: create-file build transfer-binary