.PHONY: build

build:
	go build -buildmode=plugin -o script.so script.go

run:
	docker-compose rm -f lxbot
	docker-compose rm -f script
	docker volume rm -f $(shell basename $(CURDIR))_script
	docker-compose up --build