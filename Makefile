# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gmpc android ios gmpc-cross evm all test clean
.PHONY: gmpc-linux gmpc-linux-386 gmpc-linux-amd64 gmpc-linux-mips64 gmpc-linux-mips64le
.PHONY: gmpc-linux-arm gmpc-linux-arm-5 gmpc-linux-arm-6 gmpc-linux-arm-7 gmpc-linux-arm64
.PHONY: gmpc-darwin gmpc-darwin-386 gmpc-darwin-amd64
.PHONY: gmpc-windows gmpc-windows-386 gmpc-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

gmpc:
	$(GORUN) build/ci.go install ./cmd/gmpc
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gmpc\" to launch gmpc."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gmpc.aar\" to use the library."
	@echo "Import \"$(GOBIN)/gmpc-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"
	
ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gmpc.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gmpc-cross: gmpc-linux gmpc-darwin gmpc-windows gmpc-android gmpc-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-*

gmpc-linux: gmpc-linux-386 gmpc-linux-amd64 gmpc-linux-arm gmpc-linux-mips64 gmpc-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-*

gmpc-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gmpc
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep 386

gmpc-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gmpc
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep amd64

gmpc-linux-arm: gmpc-linux-arm-5 gmpc-linux-arm-6 gmpc-linux-arm-7 gmpc-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep arm

gmpc-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gmpc
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep arm-5

gmpc-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gmpc
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep arm-6

gmpc-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gmpc
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep arm-7

gmpc-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gmpc
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep arm64

gmpc-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gmpc
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep mips

gmpc-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gmpc
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep mipsle

gmpc-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gmpc
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep mips64

gmpc-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gmpc
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-linux-* | grep mips64le

gmpc-darwin: gmpc-darwin-386 gmpc-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-darwin-*

gmpc-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gmpc
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-darwin-* | grep 386

gmpc-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gmpc
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-darwin-* | grep amd64

gmpc-windows: gmpc-windows-386 gmpc-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-windows-*

gmpc-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gmpc
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-windows-* | grep 386

gmpc-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gmpc
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gmpc-windows-* | grep amd64
