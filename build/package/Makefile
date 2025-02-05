export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

# LDFLAGS := -s -w
LDFLAGS := $${LDGOFLAGS}
BINARY := $${ProjectName}

ifeq ($(OS),Windows_NT)
    BINARY := $${ProjectName}.exe
else
    BINARY := $${ProjectName}
endif

all:  build

fmt:
	go fmt ../../...

fmt-more:
	gofumpt -l -w .

vet:
	go fmt ../../...

build: server

server:	
	source ./environment.sh;\
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)"  -o ../bin/"$(BINARY)" ../../;
clean:
	 rm -rfv ../bin/"$${ProjectName}"