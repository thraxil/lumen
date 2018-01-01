package main

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

type s3Backend struct {
	AccessKey  string
	SecretKey  string
	BucketName string
	bucket     *s3.Bucket
}

func newS3Backend(accessKey, secretKey, bucket string) *s3Backend {
	auth := aws.Auth{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	// TODO: allow configuration of buckets in other regions
	useast := aws.USEast

	connection := s3.New(auth, useast)
	mybucket := connection.Bucket(bucket)

	return &s3Backend{accessKey, secretKey, bucket, mybucket}
}

func (s s3Backend) String() string {
	return "S3"
}

func (s *s3Backend) Write(key key, r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("error writing into buffer")
		log.Println(err)
		return err
	}

	// TODO support mimetypes
	err = s.bucket.Put(key.String(), b, "application/octet", s3.BucketOwnerFull)
	if err != nil {
		log.Println("uh oh. couldn't write to bucket")
		log.Println(err)
		return err
	}
	return nil
}

func (s s3Backend) Read(key key) ([]byte, error) {
	return s.bucket.Get(key.String())
}

func (s s3Backend) Exists(key key) bool {
	ls, err := s.bucket.List(key.String(), "", "", 1)
	if err != nil {
		return false
	}
	return len(ls.Contents) == 1
}

func (s *s3Backend) Delete(key key) error {
	return s.bucket.Del(key.String())
}

func (s s3Backend) FreeSpace() uint64 {
	// TODO: this is just dummied out for now
	return 1000000000
}
