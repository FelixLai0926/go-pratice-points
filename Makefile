ENV ?= example

ifneq ("$(wildcard .env.$(ENV))","")
-include ./configs/.env.$(ENV)
endif

gen-models:
	@echo "Generating GORM models..."
	go run ./cmd/codegen/main.go
	@echo "Models generated in ./pkg/models directory."

clean-models:
	@echo "Cleaning generated models..."
	rm -rf ./pkg/models/*
	@echo "Models cleaned."
