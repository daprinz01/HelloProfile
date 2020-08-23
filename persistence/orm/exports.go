package orm

import "authengine/persistence/orm/authdb"

type Env struct {
	AuthDb *authdb.Queries
}
