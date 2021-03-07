MODULE ?= market
.PHONY:clean
all: linux

linux:
	mkdir -p dist/linux
	CGO_ENABLED=1 go build -o dist/linux/$(MODULE) main.go
    ifdef BEANPATH
	    cp dist/linux/$(MODULE) $(BEANPATH)
    endif
# arm:
# 	mkdir -p dist/arm
# 	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -mod=vendor -o dist/arm/$(MODULE) plugin/main.go

clean:
	rm -rf dist