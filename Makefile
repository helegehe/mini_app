.PHONY: build clean

DIR_SRC=./cmd/mini_app
DST=./target/bin
PACK_DIR=package
VERSION=$(shell date +'%Y.%m.%d.%H%M%S' )
PACKAGE_NAME:=mini_app_$(VERSION).tgz
GO=$(shell which go)

clean:
	@rm -rf ${PACK_DIR}
	@rm -rf $(PACKAGE_NAME)
	@rm -rf ${DST}

build: clean
	@mkdir -p ${DST}/linux-amd64
	@GOOS=linux GOARCH=amd64 $(GO) build -o ${DST}/linux-amd64 ${DIR_SRC}

pack: build
	@mkdir -p ${PACK_DIR}/bin
	@cp  ${DST}/linux-amd64/* ${PACK_DIR}/bin
	@tar -czvf $(PACKAGE_NAME) -C ${PACK_DIR} .
	@rm -rf ${PACK_DIR}
