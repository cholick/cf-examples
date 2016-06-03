package kv

type memoryStore struct {
	data map[string](map[string][]byte)
}

func NewMemoryStore() Store {
	return &memoryStore{
		data: make(map[string](map[string][]byte)),
	}
}

func (this memoryStore) CreateBucket(bucketName string) {
	bucket := this.data[bucketName]
	if bucket == nil {
		bucket = make(map[string][]byte)
		this.data[bucketName] = bucket
	}
}

func (this memoryStore) Get(bucketName string, key string) ([]byte, error) {
	if !this.bucketExists(bucketName) {
		return nil, newBucketDoesNotExistError(bucketName)
	}

	bucket := this.data[bucketName]
	return bucket[key], nil
}

func (this memoryStore) List(bucketName string) ([]string, error) {
	if !this.bucketExists(bucketName) {
		return nil, newBucketDoesNotExistError(bucketName)
	}

	bucket := this.data[bucketName]

	keys := make([]string, 0, len(bucket))
	for key := range bucket {
		keys = append(keys, key)
	}
	return keys, nil
}

func (this memoryStore) Set(bucketName string, key string, val []byte) error {
	if !this.bucketExists(bucketName) {
		this.CreateBucket(bucketName)
	}

	bucket := this.data[bucketName]
	bucket[key] = val
	return nil
}

func (this memoryStore) Del(bucketName string, key string) error {
	if !this.bucketExists(bucketName) {
		return newBucketDoesNotExistError(bucketName)
	}

	bucket := this.data[bucketName]
	delete(bucket, key)
	return nil
}

func (this memoryStore) bucketExists(bucketName string) bool {
	return this.data[bucketName] != nil
}

func newBucketDoesNotExistError(bucketName string) error {
	return BucketDoesNotExistError{
		bucketName: bucketName,
	}
}
