package upload

import (
	"fmt"
	"log"
	"os"
	"saverbate/pkg/broadcast"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-redsync/redsync/v4"
	"github.com/spf13/viper"
)

const bucket = "saverbate-vod"

// Upload is uploading of objects (video and pictures)
type Upload struct {
	record *broadcast.Record
	mutex  *redsync.Mutex
}

// Run runs uploading to Object storage
func (t *Upload) Run() {
	if err := t.mutex.Lock(); err != nil {
		log.Printf("ERROR: Upload of "+t.record.BroadcasterName+" already run: %v", err)
		return
	}

	defer func() {
		if ok, err := t.mutex.Unlock(); !ok || err != nil {
			log.Printf("ERROR: Could not release lock: %v", err)
			return
		}
	}()

	sess := session.Must(session.NewSession(&aws.Config{
		MaxRetries:  aws.Int(3),
		Credentials: credentials.NewEnvCredentials(),
		Endpoint:    aws.String(viper.GetString("objectStorageEndpoint")),
		Region:      aws.String(viper.GetString("objectStorageRegion")),
	}))

	uploader := s3manager.NewUploader(sess)

	location := t.record.BroadcasterName + "/" + t.record.UUID + ".%s"
	filename := "/app/downloads/" + location

	// Upload mp4
	fmp4, err := os.Open(fmt.Sprintf(filename, "mp4"))
	if err != nil {
		log.Printf("ERROR: failed to open file %q, %v", fmt.Sprintf(filename, "mp4"), err)
		return
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf(location, "mp4")),
		Body:   fmp4,
	})
	if err != nil {
		log.Printf("ERROR: failed to upload file, %v", err)
		fmp4.Close()

		return
	}
	fmp4.Close()

	log.Printf("File uploaded to, %s\n", result.Location)

	// Upload jpg
	fjpg, err := os.Open(fmt.Sprintf(filename, "jpg"))
	if err != nil {
		log.Printf("ERROR: failed to open file %q, %v", fmt.Sprintf(filename, "jpg"), err)
		return
	}
	// Upload the file to S3.
	result, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf(location, "jpg")),
		Body:   fjpg,
	})
	if err != nil {
		log.Printf("ERROR: failed to upload file, %v", err)
		fjpg.Close()
		return
	}
	log.Printf("File uploaded to, %s\n", result.Location)
	fjpg.Close()

	// Cleanup
	if err = os.Remove(fmt.Sprintf(filename, "jpg")); err != nil {
		log.Printf("ERROR: failed to remove file jpg, %v", err)
	}
	if err = os.Remove(fmt.Sprintf(filename, "mp4")); err != nil {
		log.Printf("ERROR: failed to remove file mp4, %v", err)
	}
}
