package controllers

import "helloprofile.com/persistence/orm/helloprofiledb"

// Env is used to declare public variable accessible to the controllers
type Env struct {
	HelloProfileDb *helloprofiledb.Queries
}
