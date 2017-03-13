// Package s3 provides a client for AWS S3.
//
package s3

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awsutil"
  "github.com/aws/aws-sdk-go/service/s3"
  "fmt"
  "os"
  "bytes"
  "net/http"
  "strings"
)

func (c *S3) UploadSingleFile(bucket_path string, file_name string, versionLabel string){
  file, err := os.Open(file_name)
  if err != nil {
    fmt.Printf("err opening file: %s", err)
  }
  defer file.Close()


  fileInfo, _ := file.Stat()
  var size int64 = fileInfo.Size()

  buffer := make([]byte, size)
  file.Read(buffer)
  fileBytes := bytes.NewReader(buffer)
  fileType := http.DetectContentType(buffer)

  bucket_info := c.ParseS3Bucket(bucket_path)

  fmt.Printf("Uploading [ %s ] to S3 bucket : %s", file.Name(), bucket_path)

  //fmt.Println(bucket_info[1] + "/" + versionLabel + ".zip")

  params := &s3.PutObjectInput{
    Bucket: aws.String(bucket_info[0]),
    Key: aws.String(bucket_info[1] + "/" + versionLabel + ".zip"),
    Body: fileBytes,
    ContentLength: aws.Int64(size),
    ContentType: aws.String(fileType),
  }

  resp, err := c.Service.PutObject(params)
  if err != nil {
    fmt.Printf("bad response: %s", err)
  }
  fmt.Printf("response %s", awsutil.StringValue(resp))
}

func (c *S3) ParseS3Bucket(bucket_path string) [2]string{
  bucket_path_split := strings.Split(bucket_path, "/")
  var rtn [2]string

  rtn[0] = bucket_path_split[0]

  rtn[1] = strings.Join(bucket_path_split[1:],"/")

  return rtn
}
