package main

import (
	"github.com/greenac/s3upload/file"
	"github.com/joho/godotenv"
	"github.com/greenac/s3upload/logger"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	ul := os.Getenv("USE_LOCAL")
	if ul == "true" {
		bp := os.Getenv("BASE_PATH")
		tp := os.Getenv("TARGET_PATH")
		if bp == "" || tp == "" {
			logger.Error("BASE_PATH and/or TARGET_PATH must be set in .env file")
			panic("Missing Env Variable")
		}

		file.GetFilesLocal(bp, tp)
	} else {
		bp := os.Getenv("BASE_PATH")
		bkp := os.Getenv("BUCKET")
		if bp == "" || bkp == "" {
			logger.Error("BASE_PATH and/or BUCKET must be set in .env file")
			panic("Missing Env Variable")
		}

		file.GetFiles(bp, bkp)
	}
}
