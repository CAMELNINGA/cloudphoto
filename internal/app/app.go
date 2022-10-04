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
	init := true
	client := client.NewAdapter()
	if err = client.InitClient(cfg); err != nil {
		init = false
	}
	service := domain.NewService(client, init, cfg)
	err = cli.NewAdapter(service)
	if err != nil {
		panic(err)
	}
}
