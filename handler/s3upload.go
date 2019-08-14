package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hypeman/metadata"
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
	S3_BUCKET = "hypemans3"
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

//Uploads video to s3 bucket and also adds metadata to datastore.
func (h *Handler) S3UploadHandler(w http.ResponseWriter, r *http.Request) {
	h.enableCors(&w)

	file, header, err := r.FormFile("video")

	if err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	username := r.FormValue("username")

	tags := []string{r.FormValue("tag1"), r.FormValue("tag2"), r.FormValue("tag3")}

	filePath := "./tmp/" + header.Filename

	//2 Concurrent Threads
	var metadata *metadata.Metadata

	done := make(chan bool)

	go func() {
		//create tmp file
		tmpFile, _ := os.Create(filePath)
		io.Copy(tmpFile, file)

		//upload file to the s3 bucket, path "[username]/filename"
		AddFileToS3(sess, header.Filename, username)
		//delete tmp file
		os.Remove(filePath)
		done <- true
	}()

	go func() {
		//add that data to the database!
		metadata, _ = h.DataStore.UploadMetadataToDB(username, header.Filename, tags)
		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}

	fmt.Printf("\n%+v\n", metadata)

	json.NewEncoder(w).Encode(metadata)
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
