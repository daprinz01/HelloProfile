package controllers

import "persianblack.com/authengine/persistence/orm/authdb"

// Env is used to declare public variable accessible to the controllers
type Env struct {
	AuthDb *authdb.Queries
}
