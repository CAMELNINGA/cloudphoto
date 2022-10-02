package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/CAMELNINGA/cloudphoto/internal/domain"
	"github.com/urfave/cli/v2"
)

type adapter struct {
	app     *cli.App
	service *domain.Service
}

func NewAdapter(service *domain.Service) error {
	a := &adapter{
		service: service,
	}
	a.app = &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	if err := a.app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	return nil
}
