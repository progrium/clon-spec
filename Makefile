build:
	@go build

test: build
	@./clon name=John age:=28 lang:='["english","zulu"]' married:=true wife:='{"name":"Jane","age":30}'