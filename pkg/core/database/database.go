package database

type Database interface {
	Users() UserStore
	Sessions() SessionStore
}
