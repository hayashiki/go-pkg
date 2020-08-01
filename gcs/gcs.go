package gcs

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"sort"

	"cloud.google.com/go/storage"
)

// Storage Storage service interface
//go:generate mockgen -source gcs.go -destination mock_gcs/mock_gcs.go
type Storage interface {
	Put(ctx context.Context, objName string, data []byte) error
	Get(ctx context.Context, objName string) ([]byte, error)
	List(ctx context.Context, filePrefix string) ([]string, error)
	Delete(ctx context.Context, objName string) error
	MakeObjectPublic(ctx context.Context, objName string) error
}

type client struct {
	gcsClient *storage.Client
	bucket    string
}

func (c *client) Put(ctx context.Context, objName string, data []byte) error {
	w := c.gcsClient.Bucket(c.bucket).Object(objName).NewWriter(ctx)
	defer w.Close()

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

// List Fetch Multi Object name request to google cloud storage.
func (c *client) List(ctx context.Context, filePrefix string) ([]string, error) {
	it := c.gcsClient.Bucket(c.bucket).Objects(ctx, &storage.Query{Prefix: filePrefix})

	var files []string
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, objAttrs.Name)
	}

	sort.Strings(files)

	return files, nil
}

func (c *client) Delete(ctx context.Context, objName string) error {
	o := c.gcsClient.Bucket(c.bucket).Object(objName)
	if err := o.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (c *client) MakeObjectPublic(ctx context.Context, objName string) error {
	acl := c.gcsClient.Bucket(c.bucket).Object(objName).ACL()

	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}

	return nil

}

func NewGCSClient(bucket string) (Storage, error) {
	ctx := context.Background()

	gcsClient, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	return &client{
		gcsClient: gcsClient,
		bucket:    bucket,
	}, nil
}

// Get Get request to google cloud storage.
func (c *client) Get(ctx context.Context, objName string) ([]byte, error) {
	r, err := c.gcsClient.Bucket(c.bucket).Object(objName).NewReader(ctx)
	defer r.Close()

	b, err := ioutil.ReadAll(r)

	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

// URL gcs object path
func (c *client) URL(obj string) string {
	return fmt.Sprintf("https://%s/%s/%s", "storage.googleapis.com", c.bucket, obj)
}
