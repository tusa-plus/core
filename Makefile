lint:
	cd facebook && golangci-lint run .
	cd tokenstorage && golangci-lint run .
	cd utils && golangci-lint run .

test:
	cd facebook && gotest -v
	cd tokenstorage && gotest -v
	cd utils && gotest -v
