all: program

cloc:
	cloc --exclude-lang="HTML,JSON" --exclude-dir="node_modules" .

test:
	grc go test -coverprofile=cover.txt ./...
	go tool cover -html=cover.txt -o cover.html

docker:
	buildah bud -f Dockerfile -t yaasbackend

program:
	cd server
	CGO_ENABLED=0 go build
