ENV ?= development

ifneq ("$(wildcard $(ENV).env)","")
-include ./configs/$(ENV).env
endif

gen-models:
	@echo "Generating GORM models with environment $(ENV)..."
	go run ./cmd/genmodel/main.go -env=$(ENV)
	@echo "Models generated."
