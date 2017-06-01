.PHONY: test

all:
	go get github.com/mattn/gom
	gom install

test:
	gom test ./lrucache -v
	gom test ./fifocache -v
