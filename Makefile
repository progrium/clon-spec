build:
	@go build

test: build deps
	@basht tests/*
	
deps:
	@type basht || go install github.com/progrium/basht@latest