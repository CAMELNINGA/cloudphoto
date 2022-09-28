package client

import (
	"context"
	"fmt"
	"log"

	localcfg "github.com/CAMELNINGA/cloudphoto/config"
	"github.com/CAMELNINGA/cloudphoto/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
)

type adapter struct {
	logger *zerolog.Logger
	config *localcfg.Config
	client *s3.Client
}

func NewAdapter(logger *zerolog.Logger, config *localcfg.Config) (domain.Client, error) {
	a := &adapter{
		logger: logger,
		config: config,
	}
	a.client = a.createClinet()
	return a, nil
}

func (a *adapter) createClinet() *s3.Client {
	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == "ru-central1" {
			return aws.Endpoint{
				PartitionID:   "yc",
				URL:           a.config.EndpointUrl,
				SigningRegion: a.config.Region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	// Подгружаем конфигрурацию из ~/.aws/*
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatal(err)
	}

	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)
	return client
}

func (a *adapter) LoadDefaultConfig() {

	// Запрашиваем список бакетов
	result, err := a.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range result.Buckets {
		log.Printf("backet=%s creation time=%s", aws.ToString(bucket.Name), bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"))
	}

}

func (a *adapter) ListObject() {
	// Запрашиваем список бакетов
	object, err := a.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(a.config.Bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range object.Contents {
		log.Printf("object=%s size=%d Bytes last modified=%s", aws.ToString(object.Key), object.Size, object.LastModified.Format("2006-01-02 15:04:05 Monday"))
	}
}
