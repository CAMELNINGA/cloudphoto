package app

import (
	"github.com/CAMELNINGA/cloudphoto/config"
	"github.com/CAMELNINGA/cloudphoto/internal/usecase/client"
)

func Run() {
	cfg, err := config.NewConfig()
	if err != nil {

	}
	//sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	client := client.NewAdapter()
	client.InitClient(cfg)
}
