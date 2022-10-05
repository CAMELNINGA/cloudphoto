package domain

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/CAMELNINGA/cloudphoto/config"
)

type Service interface {
	InitClient(bucket, awskey, awssec string) error
	Download(album, dir string) error
	Upload(album, path string) error
	List(album string) error
	Delete(album, photo string) error
	MkSite() (string, error)
}

type service struct {
	client Client
	init   bool
	config *config.Config
}

func NewService(client Client, init bool, config *config.Config) Service {
	return &service{
		init:   init,
		client: client,
		config: config,
	}
}

func (s *service) InitClient(bucket, awskey, awssec string) error {
	config, err := config.InitConfig(bucket, awskey, awssec)
	if err != nil {
		return err
	}
	s.init = true
	s.config = config
	return s.client.InitClient(config)
}

func (s *service) Download(album, dir string) error {
	if !s.init {
		return fmt.Errorf("init config pls")
	}
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
		if err := os.MkdirAll(dir, 0770); err != nil {
			fmt.Printf("Error creating dir %s \n", filepath.Dir(object))
			break
		}
		wg.Add(1)

		go d(object, dir+"/"+a[1])
	}

	wg.Wait()
	return nil
}

func (s *service) Upload(album, path string) error {
	if !s.init {
		return fmt.Errorf("init config pls")
	}
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
		p := regexp.MustCompile(`.`)
		filea := a.Split(file, 2)
		filep := p.Split(file, -1)
		ty := len(filep)
		fileType := make([]byte, 512)
		f.Read(fileType)
		types := http.DetectContentType(fileType)
		if types != "image/jpeg" && filep[ty-1] != "jpeg" && filep[ty-1] != "jpg" {
			fmt.Println(types)
			return fmt.Errorf("Invaid files")
		}
		if err := s.client.PutObject(f, album+"/"+filea[1], fi.Size(), types); err != nil {
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
	a := regexp.MustCompile(`/`)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			splitPath := a.Split(path, 3)
			if len(splitPath) == 2 {
				files = append(files, path)
			}
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
	albums := make(map[string]bool, 0)
	if !s.init {
		return fmt.Errorf("init config pls")
	}
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
		if album == "" {
			albums[a[0]] = true
		} else if len(a) != 1 {
			fmt.Printf("Object %s \n", a[1])
		}
	}
	for i := range albums {
		fmt.Println(i)
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
	if !s.init {
		return fmt.Errorf("init config pls")
	}
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

type Url struct {
	Url  string
	Name string
}
type Body struct {
	Urls  []Url
	Index string
}

func (s *service) albumfHtml(data *Body) (*bytes.Buffer, error) {
	check := func(b *bytes.Buffer, err error) (*bytes.Buffer, error) {
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	var tpl bytes.Buffer
	t, err := template.New("webpage").Parse(albumHtml)
	if err != nil {
		return nil, err
	}
	err = t.Execute(&tpl, data)

	return check(&tpl, err)
}

func (s *service) indexfHtml(data *Body) (*bytes.Buffer, error) {
	check := func(b *bytes.Buffer, err error) (*bytes.Buffer, error) {
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	var tpl bytes.Buffer
	t, err := template.New("webpage").Parse(indexHtml)
	if err != nil {
		return nil, err
	}
	err = t.Execute(&tpl, data)

	return check(&tpl, err)
}

func (s *service) errorfHtml(data *Url) (*bytes.Buffer, error) {
	check := func(b *bytes.Buffer, err error) (*bytes.Buffer, error) {
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	var tpl bytes.Buffer
	t, err := template.New("webpage").Parse(errorHtml)
	if err != nil {
		return nil, err
	}
	err = t.Execute(&tpl, data)

	return check(&tpl, err)
}

func (s *service) MkSite() (string, error) {
	albums := make(map[string]bool, 0)
	objects, err := s.client.ListObject()
	if err != nil {
		fmt.Printf("Error while getig objects %s \n", err)
		return "", err
	}
	a := regexp.MustCompile(`/`)
	for _, object := range objects {
		a := a.Split(object, 2)
		albums[a[0]] = true
	}
	baseurl := s.config.EndpointUrl + "/" + s.config.Bucket
	indexU := []Url{}
	ii := 0
	for i := range albums {

		u := []Url{}

		for _, object := range objects {
			a := a.Split(object, 2)
			if i == a[0] && len(a) == 2 {
				u = append(u, Url{
					Url:  baseurl + "/" + a[0] + "/" + a[1],
					Name: a[1],
				})
			}
		}
		data := Body{
			Urls:  u,
			Index: baseurl + "/" + "index.html",
		}
		b, err := s.albumfHtml(&data)
		if err != nil {
			return "", fmt.Errorf("Error while creating html %s \n", i)
		}
		ib := b.Bytes()
		r := bytes.NewReader(ib)
		types := http.DetectContentType(ib)
		if err := s.client.PutObject(r, "album"+strconv.Itoa(ii)+".html", 0, types); err != nil {
			return "", fmt.Errorf("Error while creating html %s \n", "album"+strconv.Itoa(ii))
		}
		if len(u) != 0 {
			indexU = append(indexU, Url{
				Url:  baseurl + "/" + "album" + strconv.Itoa(ii) + ".html",
				Name: i,
			})
			ii++
		}
	}
	data := Body{
		Urls: indexU,
	}
	index, err := s.indexfHtml(&data)

	if err != nil {
		return "", fmt.Errorf("Error while creating html %s \n", "index")
	}
	ib := index.Bytes()
	r := bytes.NewReader(ib)

	types := http.DetectContentType(ib)

	if err := s.client.PutObject(r, "index.html", 0, types); err != nil {
		return "", fmt.Errorf("Error while creating html %s \n", "index")
	}
	u := Url{
		Url:  baseurl + "/" + "index.html",
		Name: "index.html",
	}
	errorh, err := s.errorfHtml(&u)
	if err != nil {
		return "", fmt.Errorf("Error while creating html %s \n", "error")
	}
	eb := errorh.Bytes()
	r = bytes.NewReader(eb)
	types = http.DetectContentType(eb)
	if err := s.client.PutObject(r, "error.html", 0, types); err != nil {
		return "", fmt.Errorf("Error while creating html %s \n", "error")
	}
	return baseurl + "/" + "index.html", nil
}
