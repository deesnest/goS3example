/*
   This is a simple code example for connecting, uploading, downloading and listing files
   from an AWS S3 Bucket.
   Author: Antonio Sanchez antonio@asanchez.dev
*/

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	AWS_S3_REGION = "ap-south-1"
	AWS_S3_BUCKET = "golangbucket"
)

var sess = connectAWS()

func connectAWS() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		panic(err)
	}
	return sess
}

func main() {

	http.HandleFunc("/upload/", handlerUpload) // Upload
	http.HandleFunc("/get/", handlerDownload)  // Get the file
	http.HandleFunc("/list/", handlerList)     // List all files
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func showError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, message)
}

func handlerUpload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	// Get a file from the form input name "file"
	file, header, err := r.FormFile("file")
	if err != nil {
		showError(w, r, http.StatusInternalServerError, "Something went wrong retrieving the file from the form")
		return
	}
	defer file.Close()

	filename := header.Filename

	// Upload the file to S3.
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET), // Bucket
		Key:    aws.String(filename),      // Name of the file to be saved
		Body:   file,                      // File
	})
	if err != nil {
		// Do your error handling here
		showError(w, r, http.StatusInternalServerError, "Something went wrong uploading the file")
		return
	}

	fmt.Fprintf(w, "Successfully uploaded to %q\n", AWS_S3_BUCKET)
	return
}
