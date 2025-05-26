package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	postmodel "github.com/JaxonAdams/blog-backend/src/models/posts"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Service struct {
	client        *s3.Client
	presignClient *s3.PresignClient
}

func New(ctx context.Context) *S3Service {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)
	return &S3Service{
		client:        client,
		presignClient: presigner,
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

func (s S3Service) GetPostHtmlURL(post postmodel.Post, ctx context.Context) (string, error) {
	bucket := os.Getenv("S3_BUCKET_NAME")
	expirySeconds, _ := strconv.Atoi(os.Getenv("S3_URL_EXPIRY_SECONDS"))
	expiry := time.Duration(expirySeconds) * time.Second

	return s.getPresignedGetURL(bucket, post.HtmlS3Key, expiry, ctx)
}

func (s S3Service) GetPostMdURL(post postmodel.Post, ctx context.Context) (string, error) {
	bucket := os.Getenv("S3_BUCKET_NAME")
	expirySeconds, _ := strconv.Atoi(os.Getenv("S3_URL_EXPIRY_SECONDS"))
	expiry := time.Duration(expirySeconds) * time.Second

	return s.getPresignedGetURL(bucket, post.MdS3Key, expiry, ctx)
}

func (s S3Service) getPresignedGetURL(bucket, key string, expiry time.Duration, ctx context.Context) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	request, err := s.presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
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
