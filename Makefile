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
dockerRun:
	 docker run --mount source=persian-black-log,destination=/usr/local/bin/log/ -p 8083:8083 --name authengine authengine:latest
.PHONY: migrationCreate migrationUp migrationDown migrationForce migrationGoto installSqlc initialiseGoModules dockerRun