// Command cloudfront-autoindex is an AWS Lambda that processes S3 events and
// makes object copies with /index.html suffix removed in their keys.
//
// The use case is to for hosting static sites on S3 fronted by CloudFront, and
// to have a similar to what some http servers provide by automatically serving
// `dir/index.html` when client requests URLs like https://example.org/dir or
// https://example.org/dir/
package main

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() { lambda.Start(handler) }

func handler(ctx context.Context, s3Event events.S3Event) error {
	const suffix = "/index.html"
	var objects []objectMeta
	for _, record := range s3Event.Records {
		if !strings.HasPrefix(record.EventName, "ObjectCreated:") {
			continue
		}
		var err error
		bucket, key := record.S3.Bucket.Name, record.S3.Object.Key
		if key, err = url.PathUnescape(key); err != nil {
			return err
		}
		if key == suffix || key == suffix[1:] || !strings.HasSuffix(key, suffix) {
			continue
		}
		objects = append(objects, objectMeta{bucket: bucket, key: key})
	}
	if len(objects) == 0 {
		return nil
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	svc := s3.NewFromConfig(cfg)
	for _, obj := range objects {
		s := strings.TrimSuffix(obj.key, suffix)
		for _, dst := range [...]string{s, s + "/"} {
			_, err := svc.CopyObject(ctx, &s3.CopyObjectInput{
				Bucket:            &obj.bucket,
				Key:               aws.String(dst),
				CopySource:        aws.String(url.PathEscape(path.Join(obj.bucket, obj.key))),
				ACL:               types.ObjectCannedACLPrivate,
				MetadataDirective: types.MetadataDirectiveCopy,
				TaggingDirective:  types.TaggingDirectiveCopy,
			})
			if err != nil {
				return fmt.Errorf("copying %q to %q: %w", obj.key, dst, err)
			}
		}
	}
	return nil
}

type objectMeta struct {
	bucket, key string
}
