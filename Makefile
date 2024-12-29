.PHONY: make

make:
	go build -o speedo main.go

update:
	git pull
	make
