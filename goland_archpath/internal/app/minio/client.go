package minio

import (
    "context"
    "fmt"
    "mime/multipart"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient interface {
    UploadFile(objectName string, file multipart.File, header *multipart.FileHeader) (string, error)
}

type client struct {
    minio *minio.Client
    bucket string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (MinioClient, error) {
    mc, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        return nil, err
    }

    ctx := context.Background()
    exists, _ := mc.BucketExists(ctx, bucket)
    if !exists {
        if err := mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
            return nil, err
        }
    }

    return &client{minio: mc, bucket: bucket}, nil
}

func (c *client) UploadFile(objectName string, file multipart.File, header *multipart.FileHeader) (string, error) {
    _, err := c.minio.PutObject(
        context.Background(),
        c.bucket,
        objectName,
        file,
        header.Size,
        minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")},
    )
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("http://localhost:9000/%s/%s", c.bucket, objectName), nil
}
