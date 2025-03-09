package cloudflare

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appConfig "github.com/adorufus/imgupper/config"
)

func NewR2Client(cfg appConfig.CloudflareConfig) (*s3.Client, error) {

	conf, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Cloudflare R2: %w", err)
	}

	client := s3.NewFromConfig(conf, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountId))
	})

	fmt.Printf("Init S3 success\n")

	return client, nil
}
