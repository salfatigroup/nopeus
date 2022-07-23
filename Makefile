GO_BIN_DIR:=$(GOPATH)/bin

build:
	go build -o $(GO_BIN_DIR)/nopeus ./apps/cli/

clean:
	rm -f $(GO_BIN_DIR)/nopeus
