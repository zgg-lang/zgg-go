
BUILD_TIME	:= $(shell date +%FT%T%z)
BUILD_HASH	:= $(shell git rev-parse HEAD)
FLAGS		:= -ldflags "-X main.BUILD_TIME=${BUILD_TIME} -X main.BUILD_HASH=${BUILD_HASH}"
GOBUILD 	:= go build

.PHONY: lexer parser cmd linuxamd64 darwinamd64 darwinarm64 devtools

cmd:
	$(GOBUILD) -o bin/zgg $@/*.go

linuxamd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(FLAGS) -o bin/linux_amd64/zgg cmd/*.go

darwinamd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(FLAGS) -o bin/darwin_amd64/zgg cmd/*.go

windowsamd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(FLAGS) -o bin/windows_amd64/zgg.exe cmd/*.go

darwinarm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GOBUILD) $(FLAGS) -o bin/darwin_arm64/zgg cmd/*.go

parser: lexer
	java -jar tools/antlr-4.8-complete.jar -Dlanguage=Go -no-listener -visitor parser/Zgg*.g4

lexer:
	python3 scripts/makelexer.py

devtools:
	$(GOBUILD) -o bin/devtools $@/*.go

stdgolibs: devtools
	sh scripts/makegostd.sh
