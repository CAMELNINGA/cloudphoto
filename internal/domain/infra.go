package domain

import (
	"io"

	"github.com/CAMELNINGA/cloudphoto/config"
)

type Client interface {
	InitClient(config *config.Config) error
	Getobject(key string) ([]byte, error)
	CreateBucket(name string) error
	PutObject(file io.Reader, key string, size int64, types string) error
	ListObject() ([]string, error)
	LoadDefaultConfig() ([]string, error)
	DeleteObject(name string) error
}
