package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Service struct {
	client *s3.Client
}

func New(ctx context.Context) *S3Service {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Service{
		client: client,
	}
}

func (s S3Service) UploadPostHTML(postID, content string, ctx context.Context) (string, error) {
	bucket := os.Getenv("S3_BUCKET_NAME")
	key := fmt.Sprintf("posts/%s.html", postID)
	fileType := "text/html"

	contentReader := strings.NewReader(content)

	_, err := s.uploadFile(
		&bucket,
		&key,
		&fileType,
		contentReader,
		ctx,
	)

	return key, err
}

func (s S3Service) UploadPostMd(postID, content string, ctx context.Context) (string, error) {
	bucket := os.Getenv("S3_BUCKET_NAME")
	key := fmt.Sprintf("posts/%s.md", postID)
	fileType := "text/markdown"

	contentReader := strings.NewReader(content)

	_, err := s.uploadFile(
		&bucket,
		&key,
		&fileType,
		contentReader,
		ctx,
	)

	return key, err
}

func (s S3Service) uploadFile(bucket, key, fileType *string, fileReader io.Reader, ctx context.Context) (*s3.PutObjectOutput, error) {

	return s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      bucket,
		Key:         key,
		Body:        fileReader,
		ACL:         types.ObjectCannedACLPrivate,
		ContentType: fileType,
	})
}
