package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

type ACL string

const (
	Public  ACL = "public-read"
	Private ACL = "private"
)

func (a ACL) String() string {
	return string(a)
}

type Options struct {
	Key            string
	Secret         string
	Endpoint       string
	Region         string
	Bucket         string
	URL            string
	ForcePathStyle bool
	DisableSSL     bool
}

func New(c Client, opt Options) *Interactor {
	return &Interactor{
		client:         c,
		bucket:         opt.Bucket,
		url:            opt.URL,
		forcePathStyle: opt.ForcePathStyle,
	}
}

func NewS3Client(opt Options) *s3.S3 {
	s3config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(opt.Key, opt.Secret, ""),
		Endpoint:         aws.String(opt.Endpoint),
		Region:           aws.String(opt.Region),
		DisableSSL:       aws.Bool(opt.DisableSSL),
		S3ForcePathStyle: aws.Bool(opt.ForcePathStyle),
	}
	newSession := session.New(s3config)
	return s3.New(newSession)
}

type Client interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

type Interactor struct {
	client         Client
	bucket         string
	url            string
	forcePathStyle bool
}

func (i *Interactor) Upload(file io.ReadSeeker, filepath string, acl ACL, contentType string) error {
	object := s3.PutObjectInput{
		Bucket:      aws.String(i.bucket),
		Key:         aws.String(filepath),
		Body:        file,
		ACL:         aws.String(acl.String()),
		ContentType: aws.String(contentType),
	}

	_, err := i.client.PutObject(&object)
	if err != nil {
		return fmt.Errorf("storage.upload, err: %w", err)
	}

	return nil
}

func (i *Interactor) Download(filepath string) (io.ReadCloser, *string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(filepath),
	}

	result, err := i.client.GetObject(input)

	if err != nil {
		return nil, nil, fmt.Errorf("storage.download, err: %w", err)
	}
	return result.Body, result.ContentType, nil
}

func (i *Interactor) Remove(filepath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(filepath),
	}

	_, err := i.client.DeleteObject(input)
	if err != nil {
		return fmt.Errorf("storage.remove, err: %w", err)
	}

	return nil
}

func (i *Interactor) GetFullURL(path string) string {
	return fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", i.url, i.bucket, path)
}
