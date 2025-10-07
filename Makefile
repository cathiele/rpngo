default: test buildbase

all: test buildall

test:
	go test $(shell find . -name '*_test.go' -printf '%h\n' | sort -u)

buildbase:
	./buildbase.sh

buildall: buildbase
	./buildextras.sh

  
