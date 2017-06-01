.PHONY: test

all:
	go get github.com/mattn/gom
	gom install

SUCCESS := \e[1;32m
FAILURE := \e[1;31m
RESET   := \e[m

test:
	gom test ./lrucache -v \
		| sed ''/PASS/s//$$(printf "$(SUCCESS)PASS$(RESET)")/'' \
		| sed ''/FAIL/s//$$(printf "$(FAILURE)FAIL$(RESET)")/''
	gom test ./fifocache -v \
		| sed ''/PASS/s//$$(printf "$(SUCCESS)PASS$(RESET)")/'' \
		| sed ''/FAIL/s//$$(printf "$(FAILURE)FAIL$(RESET)")/''
