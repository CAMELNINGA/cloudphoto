package app

import (
	"github.com/CAMELNINGA/cloudphoto/config"
	"github.com/CAMELNINGA/cloudphoto/internal/domain"
	"github.com/CAMELNINGA/cloudphoto/internal/usecase/client"
)

func Run() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	//sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	client := client.NewAdapter()
	client.InitClient(cfg)
	service := domain.NewService(client)
	service.List("")
	service.Upload("test", "config")
	service.List("test")
	service.Download("test", "./")
	service.Download("test", "test")
	service.Delete("test", "")
}
