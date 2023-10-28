package papigoplug

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AWSGetSession uses credentials sourced from environment variables to establish an AWS session and return it.
func AWSGetSession() (sess *session.Session, err error) {
	var unsetEnvVars []string
	region, isSet := os.LookupEnv("AWS_REGION")
	if !isSet {
		unsetEnvVars = append(unsetEnvVars, "AWS_REGION")
	}
	accessKeyID, isSet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !isSet {
		unsetEnvVars = append(unsetEnvVars, "AWS_ACCESS_KEY_ID")
	}
	secretAccessKey, isSet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !isSet {
		unsetEnvVars = append(unsetEnvVars, "AWS_SECRET_ACCESS_KEY")
	}
	if len(unsetEnvVars) > 0 {
		err = fmt.Errorf("environment variable(s) are not set: %q", unsetEnvVars)
		return
	}

	sess, err = session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	return
}

// S3Download creates a local file a blob from s3 and saves it to a local file path, using multipart concurrency.
// The file is first downloaded to a temporary location on the disk, and then is renamed/moved to the final destination.
func S3Download(sess *session.Session, bucket, key, fileName string) (err error) {
	s3url := fmt.Sprintf("s3://%q/%q", bucket, key)
	Log.Debugf("Downloading %q -> %q", s3url, fileName)

	// Download to a temporary file
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return
	}
	defer os.Remove(tmpFile.Name())
	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
	})
	_, err = downloader.Download(tmpFile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		Log.Errorf("Failed to download %q: %s", s3url, err)
		return
	}
	err = tmpFile.Sync()
	if err != nil {
		Log.Errorf("Failed to download %q: %s", s3url, err)
		return
	}

	// Rename the temporary file to the final name
	err = os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	if err != nil {
		Log.Errorf("Failed to make directory tree %q: %s", filepath.Dir(fileName), err)
		return
	}
	err = os.Rename(tmpFile.Name(), fileName)
	if err != nil {
		Log.Errorf("After downloading %q, failed to move temp file to %q: %s", s3url, fileName, err)
		return
	}

	Log.Infof("Successfully downloaded %q", s3url)
	return
}

// S3Upload creates a blob on s3 by streaming the bytes from a local file path, using multipart concurrency.
func S3Upload(sess *session.Session, bucket, key, fileName string) (err error) {
	s3url := fmt.Sprintf("s3://%q/%q", bucket, key)
	Log.Infof("Uploading %q -> %q", fileName, s3url)

	file, err := os.Open(fileName)
	if err != nil {
		Log.Errorf("Failed to open %q: %s", fileName, err)
		return
	}
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024 // 64MB per part
	})
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		Log.Errorf("Failed to upload %q -> %q: %s", fileName, s3url, err)
		return
	}

	Log.Infof("Successfully uploaded %q -> %q", fileName, s3url)
	return
}

// S3SplitBucketKey returns the bucket and key components of a s3:// URL.
func S3SplitBucketKey(s3url string) (bucket string, key string, err error) {
	r := regexp.MustCompile(`^s3://(.+?)/(.+)`)
	matches := r.FindStringSubmatch(s3url)
	if len(matches) != 3 {
		err = fmt.Errorf(`invalid s3 URL (must be like "s3://bucketName/key/for/blob.txt"). Provided: %q`, s3url)
		return
	} else {
		bucket = matches[1]
		key = matches[2]
		return
	}
}
