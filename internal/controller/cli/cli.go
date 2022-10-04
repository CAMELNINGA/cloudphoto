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
	service domain.Service
}

func NewAdapter(service domain.Service) error {
	var album, path, photo string
	a := &adapter{
		service: service,
	}
	a.app = &cli.App{
		Name: "cloudphoto",
		Commands: []*cli.Command{
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "album",
						Aliases:     []string{"a"},
						Usage:       "album usege `ALBUM`",
						Destination: &album,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "path",
						Usage:       "download use path `PATH`",
						Destination: &path,
						DefaultText: "./",
					},
				},
				Name:  "download",
				Usage: "download photo in bucket",
				Action: func(*cli.Context) error {
					return service.Download(album, path)
				},
			},
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "album",
						Aliases:     []string{"a"},
						Usage:       "album usege `ALBUM`",
						Destination: &album,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "path",
						Usage:       "download use path `PATH`",
						Destination: &path,
						DefaultText: "./",
					},
				},
				Name:  "upload",
				Usage: "upload photo to bucket",
				Action: func(*cli.Context) error {
					return service.Upload(album, path)
				},
			},
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "album",
						Aliases:     []string{"a"},
						Usage:       "album usege `ALBUM`",
						Destination: &album,
					},
				},
				Name:  "list",
				Usage: "list photo album",
				Action: func(*cli.Context) error {
					return service.List(album)
				},
			},
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "album",
						Aliases:     []string{"a"},
						Usage:       "album usege `ALBUM`",
						Destination: &album,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "photo",
						Usage:       "download use path `PATH`",
						Destination: &photo,
						DefaultText: "",
					},
				},
				Name:  "delete",
				Usage: "delete photo in bucket",
				Action: func(*cli.Context) error {
					return service.Delete(album, photo)
				},
			},
			{
				Name:  "init",
				Usage: "init config",
				Action: func(cCtx *cli.Context) error {
					var bucket, awskey, awssec string
					fmt.Print("Введите aws_access_key_id: ")
					fmt.Fscan(os.Stdin, &awskey)
					fmt.Print("Введите aws_secret_access_key: ")
					fmt.Fscan(os.Stdin, &awssec)
					fmt.Print("Введите bucket: ")
					fmt.Fscan(os.Stdin, &bucket)
					return service.InitClient(bucket, awskey, awssec)
				},
			},
		},
	}

	if err := a.app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	return nil
}
