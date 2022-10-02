package client

import (
	"context"
	"fmt"
	"io"

	localcfg "github.com/CAMELNINGA/cloudphoto/config"
	"github.com/CAMELNINGA/cloudphoto/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type adapter struct {
	//logger *zerolog.Logger
	config *localcfg.Config
	client *s3.Client
	ctx    context.Context
}

func NewAdapter() domain.Client {
	a := &adapter{}

	return a
}

func (a *adapter) InitClient(localcfg *localcfg.Config) error {
	a.config = localcfg
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
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		//config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKID", "SECRET_KEY", "TOKEN")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		//a.logger.Err(err).Msg("Error while init config")
		return domain.ErrInternalS3
	}

	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)
	a.client = client
	return nil
}

func (a *adapter) LoadDefaultConfig() ([]string, error) {

	// Запрашиваем список бакетов
	result, err := a.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		//a.logger.Err(err).Msg("Error while upload list bucket")
		return nil, domain.ErrInternalS3
	}
	buckets := make([]string, 0)
	for _, bucket := range result.Buckets {
		fmt.Printf("backet=%s creation time=%s", aws.ToString(bucket.Name), bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"))
		buckets = append(buckets, aws.ToString(bucket.Name))
	}
	return buckets, nil
}

func (a *adapter) ListObject() ([]string, error) {
	// Запрашиваем список бакетов
	object, err := a.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(a.config.Bucket),
	})
	if err != nil {
		//a.logger.Err(err).Msg("Error while upload list bucket")
		return nil, domain.ErrInternalS3
	}
	objects := make([]string, 0)
	for _, object := range object.Contents {
		fmt.Printf("object=%s size=%d Bytes last modified=%s", aws.ToString(object.Key), object.Size, object.LastModified.Format("2006-01-02 15:04:05 Monday"))
		objects = append(objects, aws.ToString(object.Key))
	}
	return objects, nil
}

func (a *adapter) PutObject(file io.Reader, key string, size int64) error {
	_, err := a.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(a.config.Bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: size,
	})
	if err != nil {
		//a.logger.Err(err).Msg("Error while upload object in bucket")
		return domain.ErrInternalS3
	}
	return nil

}

func (a *adapter) CreateBucket(name string) error {
	_, err := a.client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		ACL:    types.BucketCannedACLPublicRead,
	})

	if err != nil {
		//a.logger.Err(err).Msg("Error while create bucket")
		return domain.ErrInternalS3
	}
	return nil
}

func (a *adapter) DeleteObject(name string) error {
	// Delete a single object.
	_, err := a.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(a.config.Bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		//a.logger.Err(err).Msg("Error while delete object in bucket")
		return domain.ErrInternalS3
	}
	return nil
}

func (a *adapter) Getobject(key string) ([]byte, error) {
	resp, err := a.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:       aws.String(a.config.Bucket),
		Key:          aws.String(key),
		RequestPayer: types.RequestPayerRequester,
	})

	if err != nil {
		//a.logger.Err(err).Msg("Error while get object in bucket")
		return nil, domain.ErrInternalS3
	}

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		//a.logger.Err(err).Msg("Error while get object in bucket")
		return nil, domain.ErrInternalS3
	}
	return file, nil
}
