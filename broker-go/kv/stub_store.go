package kv

//didn't want to export this, but used in multiple test packages
type StubStore struct {
	CreateBucketImpl func(string)

	ListImpl func(string) ([]string, error)
	SetImpl  func(string, string, []byte) error
	GetImpl  func(string, string) ([]byte, error)
	DelImpl  func(string, string) error
}

func (this StubStore) CreateBucket(bucketName string) {
	if this.CreateBucketImpl != nil {
		this.CreateBucketImpl(bucketName)
	}
}

func (this StubStore) List(bucketName string) ([]string, error) {
	if this.ListImpl == nil {
		panic("No ListImpl implementation")
	}
	return this.ListImpl(bucketName)
}

func (this StubStore) Set(bucketName string, key string, val []byte) error {
	if this.SetImpl == nil {
		panic("No SetImpl implementation")
	}
	return this.SetImpl(bucketName, key, val)
}

func (this StubStore) Get(bucketName string, key string) ([]byte, error) {
	if this.GetImpl == nil {
		panic("No GetImpl implementation")
	}
	return this.GetImpl(bucketName, key)
}

func (this StubStore) Del(bucketName string, key string) error {
	if this.DelImpl == nil {
		panic("No DelImpl implementation")
	}
	return this.DelImpl(bucketName, key)
}
