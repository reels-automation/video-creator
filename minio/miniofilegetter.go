package minio

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioFileGetter struct{ 
	endpoint string
	accessKey string
	secretKey string
	useSSL bool
	client *minio.Client

}

func NewMinioFileGetter(endpoint string, accessKey string, secretKey string, useSSL bool) (*MinioFileGetter, error){
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey,""),
		Secure: useSSL,
	})

	if err != nil{
		log.Fatal(err)
	}

	return &MinioFileGetter{
		endpoint: endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		useSSL: useSSL,
		client: minioClient,
	}, nil

}

func (m *MinioFileGetter) GetFile(bucket string, objectName string, filePath string) string {
	err := m.client.FGetObject(context.Background(), bucket, objectName, filePath, minio.GetObjectOptions{})

	if err != nil {
		log.Fatalf("Error trying to get: %s \n Error: %s", objectName,err)
	}
	return filePath
}
