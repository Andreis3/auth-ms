package adapter

type Database interface {
	Instance() any
	Close()
}
