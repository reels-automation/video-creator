package minio

type FileGetter interface{
	GetFile(bucket string, objectName string , filePath string) string
}
