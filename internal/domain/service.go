package domain

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/CAMELNINGA/cloudphoto/config"
)

type Service interface {
}

type service struct {
	client Client
}

func NewService(client Client) Service {
	return &service{

		client: client,
	}
}

func (s *service) InitClient(config *config.Config) error {
	return s.client.InitClient(config)
}

func (s *service) Download(album, dir string) error {
	objects, err := s.client.ListObject()
	if err != nil {
		fmt.Printf("Error while getig objects %s \n", err)
		return err
	}
	a := regexp.MustCompile(`/`)
	var wg sync.WaitGroup
	d := func(object, dir, name string) {
		defer wg.Done()
		b, err := s.client.Getobject(object)
		if err != nil {
			fmt.Printf("Error download object %s \n", object)
		}
		f, err := os.Create(name)
		defer f.Close()
		if err != nil {
			//s.logger.Err(err).Msg("Error  while creating file")
			return
		}

		_, err = f.Write(b)
		if err != nil {
			//s.logger.Err(err).Msg("Error  while craeting file")
			return
		}
		fmt.Printf("successfully downloaded data from %s/\n to file %s\n", object, name)
	}
	for _, object := range objects {
		a := a.Split(object, 2)
		if album != a[0] {
			continue
		}
		wg.Add(1)
		go d(object, dir, a[1])
	}

	wg.Wait()
	return nil
}
