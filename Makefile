ENV ?= example

ifneq ("$(wildcard .env.$(ENV))","")
-include .env.$(ENV)
endif

gen-models:
	@echo "Generating GORM models..."
	gormt -dsn "$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?charset=utf8mb4&parseTime=True&loc=Local" -o ./models
	@echo "Models generated in ./models directory."

clean-models:
	@echo "Cleaning generated models..."
	rm -rf ./models/*
	@echo "Models cleaned."
