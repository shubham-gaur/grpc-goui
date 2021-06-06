export CGO_ENABLED=0

APP_VERSION = 0.0.1

.PHONY:
all: clean install docker pkg

.PHONY:
prepare:
	go get github.com/fullstorydev/grpcui/...

.PHONY:
install:
	go build -o assets/grpcui github.com/fullstorydev/grpcui/cmd/grpcui
	go build -o bin/goui cmd/main.go

.PHONY:
docker: clean
	docker build -t github.com/shubham-gaur.io/goui:${APP_VERSION} -f Dockerfile .

.PHONE:
pkg:
	mkdir -p target/pkg
	cp -r bin assets templates target
	docker save github.com/shubham-gaur.io/goui:${APP_VERSION} -o target/pkg/goui-${APP_VERSION}.tar
	helm package helm/goui -d target/pkg
	mkdir -p target/tar
	tar -czvf goui-bundle.tar.gz target/assets target/bin target/pkg target/templates 

test:
	docker run -itd --rm -p 8080:8080  github.com/shubham-gaur.io/goui:${APP_VERSION}

.PHONY:
clean:
	rm -rf bin/goui assets/grpcui target
	docker rmi -f github.com/shubham-gaur.io/goui:${APP_VERSION}