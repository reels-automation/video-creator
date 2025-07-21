package minio

type FileGetter interface{
	GetFile(directory string, objectName string , filePath string) string
}
