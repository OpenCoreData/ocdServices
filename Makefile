BINARY := services
DOCKERVER :=`cat VERSION`
.DEFAULT_GOAL := linux

linux:
	cd cmd/services ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)

docker:
	docker build  --tag="opencoredata/ocdservices:$(DOCKERVER)"  --file=./build/Dockerfile .

dockerlatest:
	docker build  --tag="opencoredata/ocdservices:latest"  --file=./build/Dockerfile .
