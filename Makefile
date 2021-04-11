backend:
	@echo "### --- Building Swallowtail Infrastructure Locally --- ###"
	@echo "Please wait..."
	docker-compose -f local.yml up --build
	
frontend:
	@echo "### --- Building Swallowtail frontend --- ###"
	@echo "Please wait..."
	cd ./web && yarn start

default:
	backend frontend

