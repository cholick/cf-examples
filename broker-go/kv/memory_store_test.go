package kv_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cholick/cf-examples/broker-go/kv"
)

var _ = Describe("MemoryStore", func() {
	var memoryStore kv.Store

	BeforeEach(func() {
		memoryStore = kv.NewMemoryStore()
	})

	Describe("with bucket", func() {
		BeforeEach(func() {
			memoryStore.CreateBucket("bucket")
		})

		It("returns value when present", func() {
			err := memoryStore.Set("bucket", "key", []byte("some value"))
			Expect(err).To(BeNil())

			val, err := memoryStore.Get("bucket", "key")
			Expect(err).To(BeNil())
			Expect(string(val)).To(Equal("some value"))
		})

		It("lists values", func() {
			err := memoryStore.Set("bucket", "key1", []byte("some value"))
			Expect(err).To(BeNil())
			err = memoryStore.Set("bucket", "key2", []byte("some other value"))
			Expect(err).To(BeNil())

			val, err := memoryStore.List("bucket")
			Expect(err).To(BeNil())
			Expect(len(val)).To(Equal(2))
			Expect(val).To(ContainElement("key1"))
			Expect(val).To(ContainElement("key2"))
		})

		It("replaces existing value", func() {
			err := memoryStore.Set("bucket", "key", []byte("some value"))
			Expect(err).To(BeNil())

			err = memoryStore.Set("bucket", "key", []byte("some new value"))
			Expect(err).To(BeNil())

			val, err := memoryStore.Get("bucket", "key")
			Expect(err).To(BeNil())
			Expect(string(val)).To(Equal("some new value"))
		})

		It("returns nil if value not stored", func() {
			err := memoryStore.Set("bucket", "key", []byte("value"))

			val, err := memoryStore.Get("bucket", "not_present")

			Expect(err).To(BeNil())
			Expect(val).To(BeNil())
		})

		It("delete removes key", func() {
			err := memoryStore.Set("bucket", "key", []byte("some value"))
			Expect(err).To(BeNil())

			err = memoryStore.Del("bucket", "key")
			Expect(err).To(BeNil())

			val, err := memoryStore.Get("bucket", "key")
			Expect(err).To(BeNil())
			Expect(val).To(BeNil())
		})

		It("no error deleting non-existent key", func() {
			err := memoryStore.Del("bucket", "key")
			Expect(err).To(BeNil())

			err = memoryStore.Set("bucket", "other_key", []byte("some value"))
			Expect(err).To(BeNil())

			err = memoryStore.Del("bucket", "key")
			Expect(err).To(BeNil())
		})
	})

	Describe("when bucket not present", func() {
		It("on get errors", func() {
			_, err := memoryStore.Get("nope", "key")

			Expect(err).To(Not(BeNil()))
			Expect(err).To(BeAssignableToTypeOf(kv.BucketDoesNotExistError{}))
			Expect(err.Error()).To(ContainSubstring("nope"))
		})

		It("on set createes bucket", func() {
			err := memoryStore.Set("nope", "key", []byte("some value"))

			val, err := memoryStore.Get("nope", "key")
			Expect(err).To(BeNil())
			Expect(string(val)).To(Equal("some value"))
		})

		It("on delete errors", func() {
			err := memoryStore.Del("nope", "key")

			Expect(err).To(Not(BeNil()))
			Expect(err).To(BeAssignableToTypeOf(kv.BucketDoesNotExistError{}))
			Expect(err.Error()).To(ContainSubstring("nope"))
		})

		It("on list errors", func() {
			_, err := memoryStore.List("nope")

			Expect(err).To(Not(BeNil()))
			Expect(err).To(BeAssignableToTypeOf(kv.BucketDoesNotExistError{}))
			Expect(err.Error()).To(ContainSubstring("nope"))

		})
	})
})
