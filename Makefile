
help: 	## show this help
	@grep -e "^[0-9a-zA-Z_-]*:.*##" $(MAKEFILE_LIST) | \
	sed 's/^\(.*\):.*##\(.*\)/\1\t\2/'

tidy:  ## update deps
	go mod tidy

build: ## build
	go build $(shell go list ./... | grep -v /example/)

lint:
	#go vet $(shell go list ./... | grep -v /example/)
	go vet ./...

#list:
#	@echo $(shell go list ./... | grep -v /example/)

clean: 	## clean
	-rm langext.test

test: clean ## test
	go test -v -timeout=10s $(shell go list ./... | grep -v example/)

bench: clean	## test bench
	# -benchtime sets the minimum amount of time that the benchmark function will run
	# -run=^# filter out all of the unit test functions.
	go test -v -bench=. -benchtime=10s $(shell go list ./... | grep -v /example/)

bench-mem: clean	## test bench with memory usage
	# -run=^# filter out all of the unit test functions.
	go test -v -bench=. -benchtime=10s -benchmem $(shell go list ./... | grep -v /example/)

coverage:	## test with coverage
	#go test --converage ./...
	go test -coverprofile=coverage.txt -covermode=atomic -v -count=1 -timeout=30s -parallel=4 -failfast $(shell go list ./... | grep -v /example/)


