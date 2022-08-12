GO_BIN_DIR:=$(GOPATH)/bin

build:
	go build -o $(GO_BIN_DIR)/nopeus ./apps/cli/

clean:
	rm -f $(GO_BIN_DIR)/nopeus

release:
	make build
	cd ./apps/cli
	goreleaser --rm-dist
	cd ../../

deploy-install-script:
	gsutil cp ./scripts/install.sh gs://salfatigroup-cdn/nopeus/install.sh
