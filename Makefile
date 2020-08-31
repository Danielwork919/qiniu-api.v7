test:
	go test -v -timeout 50m -coverprofile=coverage.txt `go list ./... | grep -E  'auth|conf|cdn|storage|rtc|internal'`

unittest:
	go test -v ./auth/...
	go test -v ./conf/...
	go test -v ./internal/...
