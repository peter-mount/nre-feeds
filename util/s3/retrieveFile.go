package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

func (c *S3Credentials) RetrieveFile(fname string, file *os.File) error {
	log.Printf("Retrieving %s", fname)

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

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(c.Bucket),
			Key:    aws.String(c.Path + fname),
		})
	if err != nil {
		return err
	}

	log.Println("Downloaded", file.Name(), numBytes, "bytes")

	return nil
}
