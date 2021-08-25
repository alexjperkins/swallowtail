.PHONY: default
default: build
	docker-compose -f local.yml up --build

backend:
	cd s.satoshi &&  sudo make docker &&  cd .. && \
	cd s.googlesheets &&  sudo make docker &&  cd .. && \
	cd s.binance &&  sudo make docker &&  cd .. && \
	cd s.account &&  sudo make docker &&  cd .. && \
	cd s.discord &&  sudo make docker &&  cd .. && \
	cd s.coingecko &&  sudo make docker &&  cd .. && \
	docker-compose -f local.yml --profile backend up --build

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
	cd s.googlesheets &&  sudo make docker &&  cd .. && \
	docker-compose -f local.yml up --build  swallowtail.s.googlesheets

coingecko:
	cd s.coingecko &&  sudo make docker &&  cd .. && \
	docker-compose -f local.yml up --build  swallowtail.s.coingecko

discord:
	cd s.discord &&  sudo -E make docker &&  cd .. && \
	docker-compose -f local.yml up --build swallowtail.s.discord

account:
	cd s.account &&  sudo make docker &&  cd .. && \
	docker-compose -f local.yml up --build swallowtail.s.account

binance:
	cd s.binance &&  sudo make docker &&  cd .. && \
	docker-compose -f local.yml up --build swallowtail.s.binance

test:
	go test ./... -short

test-integration: postgres_test
	go test ./... --tags=integration

protos:
	find . -type d -name s.\* -exec bash -c './bin/generate_protobufs {}' \;
	
.PHONY: ecr-login
ecr-login:
	aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 638234331039.dkr.ecr.us-east-2.amazonaws.com
