package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TODO fill these in!
const (
	S3_REGION = "us-east-1"
	S3_BUCKET = "hypemanvideos"
)

var sess *session.Session

func init() {
	var err error
	// Create a single AWS session (we can re use this if we're uploading many files)
	sess, err = session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})

	if err != nil {
		log.Fatal(err)
	}
}
func (h *Handler) s3UploadHandler(w http.ResponseWriter, r *http.Request) {
	h.enableCors(&w)

	file, header, _ := r.FormFile("video")
	username := r.FormValue("username")

	filePath := "./tmp/" + header.Filename

	tmpFile, _ := os.Create(filePath)
	io.Copy(tmpFile, file)

	AddFileToS3(sess, header.Filename, username)

	os.Remove(filePath)

	//	AddFileToS3(sess)
	w.WriteHeader(http.StatusOK)
	return
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, fileName string, username string) error {
	// Open the file for use
	file, err := os.Open("./tmp/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3_BUCKET),
		Key:                  aws.String(username + "/" + fileName),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
