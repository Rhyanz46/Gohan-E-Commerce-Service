test:
	go test ./tests/... -v

run:
	go run main.go

run_docker:
	docker build --tag code_interview .