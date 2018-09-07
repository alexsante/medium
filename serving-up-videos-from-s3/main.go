package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {

	// Target S3 Object Key for this demo
	k := "SampleVideo_1080x720_10mb.mp4"

	// AWS credentials will be pulled from it's config file located at ~/.aws/credentials
	// but you can also provide the keys in the config via it's configuration property.
	config := &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("/Users/alexsante/.aws/credentials", "default"),
	}

	// Establishes a new S3 client by creating a new session and providing the
	// configuration details needed to connect.
	sess := session.New(config)

	s3c := s3.New(sess, config)

	// Retrieves an object on S3 by it's key.
	output, err := s3c.GetObject(&s3.GetObjectInput{Bucket: aws.String("thehive-use1"), Key: aws.String(k)})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Converts the S3 output body to a buffer
	buff, _ := ioutil.ReadAll(output.Body) // output.Body is a readcloser

	// Converts the buffer into a read seeker
	reader := bytes.NewReader(buff)

	// Where all the magic happens. ServeContent will automatically detect the mimetype by inspecting the first
	// 512 bytes of content.  If it is unable to find a specific mimetype, it will fall back to application/octet-stream.
	// ServeContent will also set the range headers for seek support. Lastly, ServeContent will also set the http status
	// code to 206 Partial Content.
	http.ServeContent(w, r, k, time.Now(), reader)

	return
}
