build:
	GOOS=linux GOARCH=amd64 go build -o ./.bin/linux/amd64/ ./
	GOOS=windows GOARCH=amd64 go build -o ./.bin/windows/amd64/ ./
	GOOS=linux GOARCH=arm go build -o ./.bin/linux/arm/ ./