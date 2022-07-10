all: install

GO=go

install: jedi.go
	@$(GO) build -gcflags=all="-N -l" -ldflags='-compressdwarf=false' .

build:
	@$(GO) build .

build2:
	@$(GO) tool go2go build

gdb:
	gdb -tui ./jedi

test:
	@$(GO) test -v

test2:
	@$(GO) tool go2go test

cover:
	@$(GO) test -cover

clean: 
	rm jedi
