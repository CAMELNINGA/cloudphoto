package domain

import "io"

type Client interface {
	Getobject(key string) ([]byte, error)
	CreateBucket(name string) error
	PutObject(file io.Reader, key string, size int64) error
	ListObject() ([]string, error)
	LoadDefaultConfig() ([]string, error)
}
