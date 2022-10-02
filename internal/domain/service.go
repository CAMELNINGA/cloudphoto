package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/CAMELNINGA/cloudphoto/config"
)

type Service interface {
	InitClient(config *config.Config) error
	Download(album, dir string) error
	Upload(album, path string) error
	List(album string) error
	Delete(album, photo string) error
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
	d := func(object, name string) {
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
		if err := os.MkdirAll(filepath.Dir(object), 0770); err != nil {
			fmt.Printf("Error creating dir %s \n", filepath.Dir(object))
			break
		}
		wg.Add(1)

		go d(object, a[1])
	}

	wg.Wait()
	return nil
}

func (s *service) Upload(album, path string) error {
	if path == "" {
		path = "./"
	}
	files, err := FilePathWalkDir(path)
	if err != nil {
		fmt.Printf("Error scaning dir %s \n", err)
		return err
	}
	op := func(path, file string) error {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error open file %s \n", err)
			return err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			fmt.Printf("Error open stat %s \n", err)
			return err
		}
		a := regexp.MustCompile(`/`)
		filea := a.Split(file, 2)
		if err := s.client.PutObject(f, album+"/"+filea[1], fi.Size()); err != nil {
			fmt.Printf("Error put object %s \n", album+"/"+file)
			return err
		}
		return nil
	}

	for _, file := range files {
		if err := op(path, file); err != nil {
			continue
		}
	}
	return err
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func haveAlbum(objects []string, album string, a *regexp.Regexp) bool {
	haveAlbum := false
	for _, object := range objects {
		a := a.Split(object, 2)
		if album != a[0] && album != "" {
			continue
		}
		haveAlbum = true
	}
	return haveAlbum
}

func (s *service) List(album string) error {
	objects, err := s.client.ListObject()
	if err != nil {
		fmt.Printf("Error while getig objects %s \n", err)
		return err
	}
	a := regexp.MustCompile(`/`)
	haveAlbum := haveAlbum(objects, album, a)
	if !haveAlbum {
		fmt.Printf("Error don't find album %s \n", album)
		return fmt.Errorf("error don't find album %s \n", album)
	}
	for _, object := range objects {
		a := a.Split(object, 2)
		if album != a[0] && album != "" {
			continue
		}
		if len(a) != 1 {
			fmt.Printf("Album %s Object %s ", a[0], a[1])
		} else {
			fmt.Printf("Object %s ", object)
		}
	}
	return nil
}

func haveObject(objects []string, sobject string) bool {
	haveAlbum := false
	for _, object := range objects {

		if object != sobject {
			continue
		}
		haveAlbum = true
	}
	return haveAlbum
}
func (s *service) Delete(album, photo string) error {
	objects, err := s.client.ListObject()
	if err != nil {
		fmt.Printf("Error while getig objects %s \n", err)
		return err
	}
	var sobject string
	if photo != "" {
		sobject = album + "/" + photo
		if !haveObject(objects, sobject) {
			fmt.Printf("Error while searching objects %s \n", sobject)
			return fmt.Errorf("error while searching objects %s \n", sobject)
		}
		if err = s.client.DeleteObject(sobject); err != nil {
			fmt.Printf("Error while deleting objects %s \n", sobject)
			return fmt.Errorf("error while deleting objects %s \n", sobject)
		}
		return nil
	}
	sobject = album
	a := regexp.MustCompile(`/`)
	if !haveAlbum(objects, sobject, a) {
		fmt.Printf("Error while searching objects %s \n", sobject)
		return fmt.Errorf("error while searching objects %s \n", sobject)
	}
	for _, object := range objects {
		a := a.Split(object, 2)
		if album != a[0] && album != "" {
			continue
		}
		if err = s.client.DeleteObject(object); err != nil {
			fmt.Printf("Error while deleting objects %s \n", object)
			return fmt.Errorf("error while deleting objects %s \n", object)
		}
	}
	return nil
}
