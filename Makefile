build:
	docker build -t main .
run:
	docker-compose -f stack.yml up
all:
	build
	run