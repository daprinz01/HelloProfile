migrationCreate:
	migrate create -ext sql -dir persistence/migrations -seq init_schema
migrationUp:
	migrate -path persistence/migrations -database "postgres://xenonprimus:Sarah4Daprinz@localhost/helloprofile?sslmode=disable" -verbose up
migrationUpTest:
	migrate -path persistence/migrations -database "postgres://ezbtavev:Hoi_daLhSQ8xiStkyLcbfLVey5QHtuAO@heffalump.db.elephantsql.com/ezbtavev?sslmode=disable" -verbose up
migrationDown:
	migrate -path persistence/migrations -database "postgres://xenonprimus:Sarah4Daprinz@localhost/helloprofile?sslmode=disable" -verbose down 
migrationForce:
	migrate -path persistence/migrations -database "postgres://xenonprimus:Sarah4Daprinz@localhost/helloprofile?sslmode=disable" -verbose force 2
migrationGoto:
	migrate -path persistence/migrations -database "postgres://xenonprimus:Sarah4Daprinz@localhost/helloprofile?sslmode=disable" -verbose goto 2
installSqlc:
	go get github.com/kyleconroy/sqlc/cmd/sqlc
initialiseGoModules:
	go mod init helloprofile
sqlcGenerate:
	sqlc generate
dockerBuild:
	docker build -t helloprofile:latest .
dockerRun:
	 docker run --mount source=persian-black-logs,destination=/usr/local/bin/log/ -p 8002:8083 --name helloprofile --env TOKEN_LIFESPAN=96h --env SESSION_LIFESPAN=120h --env DB_HOST=host.docker.internal --env DB_USER=xenonprimus --env DB_NAME=helloprofile --env DB_PASSWORD=Sarah4Daprinz --env DB_PORT=5432 --env COMMUNICATION_SERVICE_ENDPOINT=http://host.docker.internal:8084 helloprofile:latest
.PHONY: migrationCreate migrationUp migrationDown migrationForce migrationGoto installSqlc initialiseGoModules dockerRun sqlcGenerate dockerBuild