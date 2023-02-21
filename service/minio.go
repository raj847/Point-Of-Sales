package service

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient() (*minio.Client, error) {
	endpoint := "is3.cloudhost.id"
	accessKeyID := "XFE8741VR80Y9MGI6MMQ"
	secretAccessKey := "cFekU0FNwOTmZdVkxtOhQGVf9kfZAAdIjPxoRzUi"
	useSSL := true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	return minioClient, err
}

func UploadToCloud(ctx context.Context, minioClient *minio.Client, filePdf *os.File, fileName string) (string, error) {
	_, err := minioClient.PutObject(ctx, "rajendra", fileName, filePdf, -1, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		ContentType: "application/pdf",
	})

	fileUploadedLink := fmt.Sprintf("https://is3.cloudhost.id/rajendra/%s", fileName)
	return fileUploadedLink, err
}
