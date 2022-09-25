run: main.go
	go run main.go

build: main.go
	go build -o main main.go
	
test: main.go
	go test ./... -v

docker-build: Dockerfile
	docker build -t slash/badidea .

docker-run: Dockerfile
	docker run -p 8080:8080 slash/badidea
