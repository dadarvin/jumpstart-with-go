export NOW=$(shell date +"%Y/%m/%d")

generate-mock:
	@echo "${NOW} == GENERATING MOCK FILES"
	@go generate .dev/generate.go

dev:
	@echo "${NOW} == Running GO server"
	@go run main.go

test:
	@echo "{NOW} == Testing Projects"
	@go test ./... -cover -race -short | tee test.out
	@.dev/test.sh test.out 15