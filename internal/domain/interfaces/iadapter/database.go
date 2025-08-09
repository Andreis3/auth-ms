package iadapter

type IDatabase interface {
	Instance() any
	Close()
}
