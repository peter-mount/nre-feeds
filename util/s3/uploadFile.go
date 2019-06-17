package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/peter-mount/golib/kernel/logger"
	"os"
)

func (c *S3Credentials) UploadFile(log *logger.Logger, file *os.File, fname string) error {
	log.Printf("Uploading %s", file.Name())

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(c.Region),
			Credentials: credentials.NewStaticCredentials(
				c.AccessKey,
				c.SecretKey,
				"",
			),
		})
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(
		&s3manager.UploadInput{
			Bucket: aws.String(c.Bucket),
			Key:    aws.String(c.Path + fname),
			Body:   file,
		})
	if err != nil {
		return err
	}

	log.Println("Uploaded", fname)

	return nil
}
