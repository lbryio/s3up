package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s <key> <secret> <region> <bucket> <path-to-file>\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	key := os.Args[1]
	secret := os.Args[2]
	region := os.Args[3]
	bucket := os.Args[4]
	filepath := os.Args[5]

	creds := credentials.NewStaticCredentials(key, secret, "")
	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds).WithMaxRetries(3)
	uploader := s3manager.NewUploader(session.Must(session.NewSession(cfg)))

	f, err := os.Open(filepath)
	checkErr(err, "failed to open file")
	defer f.Close()

	fmt.Println("uploading...")
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path.Base(filepath)),
		Body:   f,
	})
	checkErr(err, "failed to upload file")

	fmt.Printf("file uploaded to %s\n", result.Location)
}

func checkErr(e error, prefix string) {
	if e == nil {
		return
	}

	if strings.Contains(e.Error(), "SignatureDoesNotMatch") {
		fmt.Printf("\nThe following error probably means your AWS secret key (second arg) is wrong:\n\n")
	} else if strings.Contains(e.Error(), "InvalidAccessKeyId") {
		fmt.Printf("\nThe following error probably means your AWS access key (first arg) is wrong:\n\n")
	}

	if prefix != "" {
		fmt.Printf("%s: ", prefix)
	}
	fmt.Printf("%v\n", e)
	os.Exit(1)
}
