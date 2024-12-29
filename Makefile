.PHONY: build

build:
	go build -o speedo main.go

update:
	git pull
	make
