migrationCreate:
	migrate create -ext sql -dir persistence/migrations -seq init_schema
migrationUp:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/authengine?sslmode=disable" -verbose up
migrationDown:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/authengine?sslmode=disable" -verbose down 
migrationForce:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/authengine?sslmode=disable" -verbose force 1
migrationGoto:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/authengine?sslmode=disable" -verbose goto 2
installSqlc:
	go get github.com/kyleconroy/sqlc/cmd/sqlc
initialiseGoModules:
	go mod init authengine
sqlcGenerate:
	sqlc generate
dockerBuild:
	docker build -t authengine:latest .
dockerRun:
	 docker run --mount source=persian-black-logs,destination=/usr/local/bin/log/ -p 8083:8083 --name authengine --env TOKEN_LIFESPAN=96h --env SESSION_LIFESPAN=120h --env DB_HOST=host.docker.internal --env DB_PORT=8669 --env COMMUNICATION_SERVICE_ENDPOINT=http://host.docker.internal:8084 authengine:latest
.PHONY: migrationCreate migrationUp migrationDown migrationForce migrationGoto installSqlc initialiseGoModules dockerRun sqlcGenerate