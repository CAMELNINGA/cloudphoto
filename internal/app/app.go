package app

import (
	"github.com/CAMELNINGA/cloudphoto/config"
	"github.com/CAMELNINGA/cloudphoto/internal/controller/cli"
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
	err = cli.NewAdapter(service)
	if err != nil {
		panic(err)
	}
}
