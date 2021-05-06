default: (build)

cloc:
	@cloc --exclude-lang="HTML,JSON" --exclude-dir="node_modules" .

test:
	@grc go test -coverprofile=cover.txt ./...
	@go tool cover -html=cover.txt -o cover.html

docker:
	@echo 'Bulding docker...'
	@buildah bud -f Dockerfile -t yaasbackend

js:
	@echo 'Building js...'
	@cd js && ./node_modules/.bin/gulp

build:
	@echo 'Building server...'
	@cd server && CGO_ENABLED=0 go build
