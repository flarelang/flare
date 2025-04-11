.PHONY: build run rund serve format cache

# build the project
build:
	@go build -o bin/zxlang

# run the project
run: build
	@SHOW_STACK=true ./bin/zxlang run testcodes/$(file).zx --debug

# rund the project
rund: build
		@DEBUG=true SHOW_STACK=true ./bin/zxlang run testcodes/$(file).zx --debug

# start the default http server
serve: build
	@DEBUG=true ./bin/zxlang serve testcodes/$(folder) --listenAddr=:3030

# format the project's code (dev only, not working properly)
format: build
	@DEBUG=true ./bin/zxlang format testcodes/format

# clean the project's cache
cache:
	@rm -rf .zxcache
