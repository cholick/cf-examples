package kv

import "fmt"

// https://github.com/boltdb/bolt#using-keyvalue-pairs
// https://gobyexample.com/json

type Store interface {
	CreateBucket(bucketName string)

	Set(bucketName string, key string, value []byte) error
	Get(bucketName string, key string) ([]byte, error)
	List(bucketName string) ([]string, error)
	Del(bucketName string, key string) error
}

type BucketDoesNotExistError struct {
	bucketName string
}

func (this BucketDoesNotExistError) Error() string {
	return fmt.Sprintf("Bucket [%s] does not exist", this.bucketName)
}
