package adapter

type DB interface {
	Instance() any
	Close()
}
