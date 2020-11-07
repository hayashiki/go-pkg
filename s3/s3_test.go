package s3

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"reflect"
	"testing"
)

func TestACL_String(t *testing.T) {
	tests := []struct {
		name string
		a    ACL
		want string
	}{
		{
			name: "success",
			a:    Public,
			want: "public-read",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInteractor_Download(t *testing.T) {
	s3mock := &S3mock{}
	s3mockErr := &S3mock{Error: errors.New("download")}

	opt := Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
		DisableSSL:     false,
	}
	obj := &s3.GetObjectOutput{}
	type fields struct {
		client         Client
		bucket         string
		url            string
		forcePathStyle bool
	}
	type args struct {
		filepath string
	}
	var (
		tests = []struct {
			name    string
			fields  fields
			args    args
			want    io.ReadCloser
			wantErr bool
		}{
			{
				"success",
				fields{
					client:         s3mock,
					bucket:         opt.Bucket,
					url:            opt.URL,
					forcePathStyle: opt.ForcePathStyle,
				},
				args{"test.png"},
				obj.Body,
				false,
			},
			{
				"error",
				fields{s3mockErr,
					opt.Bucket,
					opt.URL,
					opt.ForcePathStyle},
				args{"test.png"},
				nil,
				true,
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				client:         tt.fields.client,
				bucket:         tt.fields.bucket,
				url:            tt.fields.url,
				forcePathStyle: tt.fields.forcePathStyle,
			}
			got, _, err := i.Download(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Download() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInteractor_Remove(t *testing.T) {
	type fields struct {
		client         Client
		bucket         string
		url            string
		forcePathStyle bool
	}
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				client:         tt.fields.client,
				bucket:         tt.fields.bucket,
				url:            tt.fields.url,
				forcePathStyle: tt.fields.forcePathStyle,
			}
			if err := i.Remove(tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInteractor_Upload(t *testing.T) {
	type fields struct {
		client         Client
		bucket         string
		url            string
		forcePathStyle bool
	}
	type args struct {
		file        io.ReadSeeker
		filepath    string
		acl         ACL
		contentType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				client:         tt.fields.client,
				bucket:         tt.fields.bucket,
				url:            tt.fields.url,
				forcePathStyle: tt.fields.forcePathStyle,
			}
			if err := i.Upload(tt.args.file, tt.args.filepath, tt.args.acl, tt.args.contentType); (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewS3Client(t *testing.T) {
	type args struct {
		opt Options
	}
	tests := []struct {
		name string
		args args
		want *s3.S3
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewS3Client(tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewS3Client() = %v, want %v", got, tt.want)
			}
		})
	}
}
