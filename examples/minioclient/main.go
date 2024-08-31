package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := "minio.traefik.lan"
	accessKeyID := "jthtpkYjPXElMUrMEEES"
	secretAccessKey := "WmK6qmMEJ9zLonD4OTotpAgg5cb769zvQLWpQ2Gk"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.TODO()
	bucket := "odd-cat"

	ok, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		panic(err)
	}
	slog.Info("BucketExists", "ok", ok, "err", err)

	obj, err := minioClient.GetObject(ctx, bucket, "collector/events/active/2024-05-24.log", minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}
	slog.Info("GetObject", "obj", obj)

	data := make([]byte, 100)
	n, err := obj.Read(data)
	slog.Info(string(data), "n", n, "err", err)

	slog.Info(fmt.Sprintf("%#v\n", minioClient)) // minioClient is now set up
}
