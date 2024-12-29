.PHONY: build

build:
	go build -o speedo cmd/speedo/main.go

update:
	git pull
	make
