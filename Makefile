.PHONY: all build run
all: build run

build:
	docker build -t main .

run:
	docker-compose -f stack.yml up

localmongo:
	docker-compose -f stack_local.yml up

stop:
	docker-compose  -f stack.yml stop
	docker compose  -f stack.yml rm

