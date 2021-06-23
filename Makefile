.PHONY: default
default: build
	docker-compose -f local.yml up --build

build:
	find . -type d -name s.\* -exec bash -c 'cd {} && sudo make docker' \;

postgres:
	docker-compose -f local.yml up --build postgres

postgres_test:
	docker-compose -f local.yml up --build postgres_test

satoshi:
	cd s.satoshi &&  sudo make docker &&  cd .. && \
		docker-compose -f local.yml up --build swallowtail.s.satoshi

googlesheets:
	docker-compose -f local.yml up --build  swallowtail.s.googlesheets

account:
	docker-compose -f local.yml up --build swallowtail.s.account

test: postgres_test
	docker-compose -f local.yml run --rm -e swallowtail.s.account go test ./...

protos:
	find . -type d -name s.\* -exec bash -c './bin/generate_protobufs {}' \;

