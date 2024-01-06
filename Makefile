
help: 	## show this help
	@grep -e "^[0-9a-zA-Z_-]*:.*##" $(MAKEFILE_LIST) | \
	sed 's/^\(.*\):.*##\(.*\)/\1\t\2/'

tidy:  ## update deps
	go mod tidy

build: ## build
	go build ./...

lint:
	go vet ./...

clean: 	## clean
	-rm langext.test

test: clean ## test
	go test -v -timeout=10s ./...

bench: clean	## test bench
	# -benchtime sets the minimum amount of time that the benchmark function will run
	# -run=^# filter out all of the unit test functions.
	go test -v -bench=. -benchtime=10s ./...

bench-mem: clean	## test bench with memory usage
	# -run=^# filter out all of the unit test functions.
	go test -v -bench=. -benchtime=10s -benchmem ./...

coverage:	## test with coverage
	#go test --converage ./...
	go test -coverprofile=coverage.txt -covermode=atomic -v -count=1 -timeout=30s -parallel=4 -failfast ./...


