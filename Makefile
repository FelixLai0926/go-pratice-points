ENV ?= development

ifneq ("$(wildcard $(ENV).env)","")
-include ./configs/$(ENV).env
endif

gen-models:
	@echo "Generating GORM models with environment $(ENV)..."
	go run ./cmd/codegen/main.go -env=$(ENV)
	@echo "Models generated in ./pkg/models directory."

clean-models:
	@echo "Cleaning generated models..."
	rm -rf ./pkg/models/*
	@echo "Models cleaned."
