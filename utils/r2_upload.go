package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadToR2 uploads a file to Cloudflare R2 and returns the public URL
func UploadToR2(data []byte, filename string, folder string) (string, string, error) {
    accessKey := os.Getenv("R2_ACCESS_KEY_ID")
    secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
    endpoint := os.Getenv("R2_ENDPOINT")
    bucket := os.Getenv("R2_BUCKET")
    region := os.Getenv("R2_REGION")
    if region == "" {
        region = "auto"
    }
    if accessKey == "" || secretKey == "" || endpoint == "" || bucket == "" {
        return "", "", fmt.Errorf("R2 configuration missing (R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_ENDPOINT, R2_BUCKET required)")
    }
    sess, err := session.NewSession(&aws.Config{
        Region:           aws.String(region),
        Endpoint:         aws.String(endpoint),
        S3ForcePathStyle: aws.Bool(false),
        Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
    })
    if err != nil {
        return "", "", err
    }
    svc := s3.New(sess)
    // Build object key
    base := filepath.Base(filename)
    objectKey := folder + "/" + base
    // Ensure unique key (append timestamp)
    objectKey = fmt.Sprintf("%s_%d%s", folder+"/"+base[:len(base)-len(filepath.Ext(base))], time.Now().UnixNano(), filepath.Ext(base))
    // Upload
    _, err = svc.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(objectKey),
        Body:   bytes.NewReader(data),
    })
    if err != nil {
        return "", "", err
    }
    // Generate presigned download URL (valid for 30 days)
    req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(objectKey),
    })
    urlStr, _ := req.Presign(30 * 24 * time.Hour)
    return urlStr, objectKey, nil
}