.PHONY: default
default: build
	docker-compose -f local.yml up --build

build:
	find . -type d -name s.\* -exec bash -c 'cd {} && sudo make docker' \;

postgres:
	docker-compose -f local.yml up --build -V postgres

postgres_test:
	docker-compose -f local.yml up --build -V postgres_test

satoshi:
	cd s.satoshi &&  sudo make docker &&  cd .. && \
		docker-compose -f local.yml up --build swallowtail.s.satoshi

googlesheets:
	docker-compose -f local.yml up --build  swallowtail.s.googlesheets

discord:
	cd s.discord &&  sudo make docker &&  cd .. && \
	docker-compose -f local.yml up --build swallowtail.s.discord

account:
	cd s.account &&  sudo make docker &&  cd .. && \
	docker-compose -f local.yml up --build swallowtail.s.account

test:
	go test ./... -short

test-integration: postgres_test
	go test ./... --tags=integration
	
protos:
	find . -type d -name s.\* -exec bash -c './bin/generate_protobufs {}' \;
