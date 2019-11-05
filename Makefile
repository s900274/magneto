export GO111MODULE=on

export CGO_ENABLED=0

all: bld

bld: magneto

magneto:
	go build -o bin/magneto ./cmd/magneto

clean:
	@rm -f init/magneto
	@rm -rf status
	@rm -f  log/*log*
	@rm -rf ./output

cleanlog:
	@rm -f log/*log*
