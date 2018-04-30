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

var blackList = map[string]int {
	"BirthdayMysteryBox.json": 0,
	"UserHasntSeenGameStatsInfoX.json": 0,
	"UserHasntSeenLgWinStatsX.json": 0,
	"UserIsMakingProgressTowardASavingsMissionAndIsCloseToGettingRewardX.json": 0,
	"UserIsMakingProgressTowardAStreakMissionAndIsCloseToGettingRewardX.json": 0,
	"LevelXIncompleteProgress2": 0,
}


func GetFiles(basePath string, bucket string) {
	fi, err := ioutil.ReadDir(basePath)
	if err != nil {
		logger.Error("Could not read files from base path:", basePath)
		return
	}

	counter := 0
	for _, fi := range fi {
		p := path.Join(basePath, fi.Name())
		f, err := os.Open(p)
		if err != nil {
			logger.Error("Could not read file:", p, "Failed with error:", err)
			continue
		}

		suc := uploadToS3(f, bucket)
		if suc {
			counter += 1
		}
	}

	logger.Log("Uploaded:", counter, "files to s3 bucket:", bucket)
}

func GetFilesLocal(basePath string, targetPath string) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		logger.Error("Could not read files from target path:", basePath, "error:", err)
		return
	}

	for _, fi := range files {
		p := path.Join(targetPath, fi.Name())
		err = os.RemoveAll(p)
		if err != nil {
			logger.Error("Could not remove file", p, "error:", err)
			continue
		}

		logger.Log("Removed file:", p)
	}

	files, err = ioutil.ReadDir(basePath)
	if err != nil {
		logger.Error("Could not read files from base path:", basePath, "error:", err)
		return
	}

	counter := 0
	for _, fi := range files {
		_, isBlkListed := blackList[fi.Name()]
		if !strings.Contains(fi.Name(), ".json") || isBlkListed {
			logger.Warn("Not copying file:", fi.Name(), "to:", targetPath)
			continue
		}

		p := path.Join(basePath, fi.Name())
		f, err := os.Open(p)
		if err != nil {
			logger.Error("Could not read file:", p, "Failed with error:", err)
			continue
		}

		d, err := ioutil.ReadAll(f)
		if err != nil {
			logger.Error("Could not read file", p, "error:", err)
			continue
		}

		tp := path.Join(targetPath, fi.Name())
		_, err = os.Create(tp)
		if err != nil {
			logger.Error("Could not create file:", tp, "error:", err)
			continue
		}

		err = ioutil.WriteFile(tp, d, 0644)
		if err != nil {
			logger.Error("Could not write file:", tp, "error:", err)
			continue
		}

		logger.Log("Copied file:", tp)
		counter += 1
	}

	logger.Log("Moved:", counter, "files to local dir:", targetPath)
}

func uploadToS3(f *os.File, bucket string) bool {
	n := getFileName(f)
	// _, hasFile := whiteList[n]

	_, isBlkListed := blackList[n]
	if !strings.Contains(n, ".json") || isBlkListed {
		logger.Warn("Not uploading file:", n, "to s3")
		return false
	}

	logger.Log("Uploading file:", n, "to s3")

	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		logger.Error("Could not start aws session to upload:", n, "error: err")
		return false
	}

	uploader := s3manager.NewUploader(s)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(n),
		Body:   f,
	})

	if err != nil {
		logger.Error("Failed to upload file:", f.Name(), "to bucket:", bucket, "error:", err)
		return false
	}

	logger.Log("Uploaded file:", f.Name(), "to bucket:", bucket, "with result:", result)
	return true
}

func getFileName(f *os.File) string {
	parts := strings.Split(f.Name(), "/")
	return parts[len(parts) - 1]
}
