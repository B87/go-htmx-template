TMP_DIR = tmp
COV_DIR = ${TMP_DIR}/.coverage
PORT = 8080
DOCKER_IMAGE = go-htmx


run-dev:
	air
run:
	${TMP_DIR}/go-htmx

test:
	mkdir -p ${COV_DIR} && go test -v -coverprofile=${COV_DIR}/cover.out ./...

build:
	npx tailwind -i 'web/static/css/tailwind.css' \
		-o 'web/static/css/styles.css'
	go build -o ${TMP_DIR}/go-htmx

docker-build:
	docker build -t ${DOCKER_IMAGE} .

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down


# Release
version-bump:
	semver -l -c semver.yaml -s

.PHONY: run-dev run test build docker-build migrate-up migrate-down version-bump
