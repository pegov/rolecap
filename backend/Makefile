test_all:
	go test -v ./... -coverprofile cover.out && go tool cover -func cover.out

test_internal:
	go test --tags=internal -v ./... -coverprofile cover.out && go tool cover -func cover.out

html:
	go tool cover -html cover.out 