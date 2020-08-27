migrationCreate:
	migrate create -ext sql -dir persistence/migrations -seq init_schema
migrationUp:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/persian_black?sslmode=disable" -verbose up
migrationDown:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/persian_black?sslmode=disable" -verbose down
migrationForce:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/persian_black?sslmode=disable" -verbose force 17
migrationGoto:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/persian_black?sslmode=disable" -verbose goto 2
installSqlc:
	go get github.com/kyleconroy/sqlc/cmd/sqlc
initialiseGoModules:
	go mod init authengine
dockerRun:
	 docker run --mount source=logs,destination=/usr/local/bin/log --add-host=localhost:127.0.0.1 authengine:latest
.PHONY: migrationCreate migrationUp migrationDown migrationForce migrationGoto installSqlc initialiseGoModules dockerRun