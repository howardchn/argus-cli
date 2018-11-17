all:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o argus-cli-darwin main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o argus-cli-linux main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o argus-cli-win main.go


