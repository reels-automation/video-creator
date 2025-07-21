package minio

import (
	"context"
	log "github.com/sirupsen/logrus"
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

func NewMinioFileGetter(endpoint string, accessKey string, secretKey string, useSSL bool) (*MinioFileGetter){
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
	}

}

func (m *MinioFileGetter) GetFile(bucket string, objectName string, filePath string) string {
	err := m.client.FGetObject(context.Background(), bucket, objectName, filePath, minio.GetObjectOptions{})

	if err != nil {
		log.Fatalf("Error trying to get: %s \n Error: %s", objectName,err)
	}
	return filePath
}

func (m *MinioFileGetter) UploadFile(bucket string, objectName string, filePath string) {
	
	info, err := m.client.FPutObject(context.Background(),bucket, objectName, filePath, minio.PutObjectOptions{})

	if err != nil{
		log.Errorf("Error trying to upload File: %s To Bucket: %s. \nError:%s", objectName,bucket,err)
	}
	log.Infof("Successfully uploaded %s of size %d \n", objectName, info.Size)
}