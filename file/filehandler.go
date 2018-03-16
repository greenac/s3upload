package file

import (
	"io/ioutil"
	"github.com/greenac/s3upload/logger"
	"path"
	"os"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)


func GetFiles(basePath string, bucket string) {
	fi, err := ioutil.ReadDir(basePath)
	if err != nil {
		logger.Error("Could not read files from base path:", basePath)
		return
	}

	for i, fi := range fi {
		p := path.Join(basePath, fi.Name())
		logger.Log("Got file #", i, p)
		f, err := os.Open(p)
		if err != nil {
			logger.Error("Could not read file:", p, "Failed with error:", err)
			continue
		}

		uploadToS3(f, bucket)
	}
}

func uploadToS3(f *os.File, bucket string) {
	parts := strings.Split(f.Name(), "/")
	n := parts[len(parts) - 1]
	if !strings.Contains(n, ".json") {
		return
	}

	//if n != "TestName.json" {
	//	return
	//}

	s, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	uploader := s3manager.NewUploader(s)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(n),
		Body:   f,
	})

	if err != nil {
		logger.Error("Failed to upload file:", f.Name(), "to bucket:", bucket, "error:", err)
		return
	}

	logger.Log("Uploaded file:", f.Name(), "to bucket:", bucket, "with result:", result)
}