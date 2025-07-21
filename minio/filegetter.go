package minio

type FileGetter interface{
	GetFile(directory string, objectName string , filePath string) string
	UploadFile(directory string, objectName string, filePath string)
}
