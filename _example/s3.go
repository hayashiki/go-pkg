package main

import (
	"github.com/hayashiki/go-pkg/s3"
	"log"
	"os"
)

func main() {
	opt := s3.Options{
		Key:            os.Getenv("AWS_ACCESS_KEY_ID"),
		Secret:         os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:         os.Getenv("AWS_REGION"),
		Bucket:         os.Getenv("AWS_BUCKET"),
		ForcePathStyle: false,
		DisableSSL:     false,
	}
	s3Service := s3.New(s3.NewS3Client(opt), opt)

	filePath := "./testdata/image.png"
	file, err := os.Open(filePath)
	if err != nil {
		log.Print("Fail to open the file")
		return
	}
	defer file.Close()

	err = s3Service.Upload(file, filePath, s3.Public, "multipart/form-data")

	if err != nil {
		log.Print("Fail to upload the file")
		return
	}
}
