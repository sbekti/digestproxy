build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o digestproxy main.go

docker:
	docker build -t sbekti/digestproxy:latest .

push:
	docker push sbekti/digestproxy:latest

run:
	go run main.go

clean:
	rm -f digestproxy

all: clean build
