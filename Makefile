build:
	GOOS=windows GOARCH=amd64 go build -o ./dist/mj.exe ./main.go
	GOOS=linux GOARCH=amd64 go build -o ./dist/mj ./main.go
	git add .
	git commit -m "build"
	git push