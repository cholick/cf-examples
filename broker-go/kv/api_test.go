package kv_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"

	"github.com/cholick/cf-examples/broker-go/kv"
)

var _ = Describe("API", func() {
	const path = "/kv"

	var _ = Describe("list", func() {
		makeRequest := func() *httptest.ResponseRecorder {
			api := kv.NewApi(lagertest.NewTestLogger("api"), &kv.StubStore{
				ListImpl: func(string) ([]string, error) {
					return []string{"key1", "key2"}, nil
				},
			})

			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, path+"/id1/", nil)

			api.ServeHTTP(recorder, request)

			return recorder
		}

		It("sets content type", func() {
			response := makeRequest()
			contentType := response.Header().Get("Content-Type")
			Expect(contentType).To(Equal("text/plain"))
		})

		It("responds with 200", func() {
			response := makeRequest()

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("responds with key list", func() {
			response := makeRequest()

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).To(BeNil())

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(string(body)).To(Equal("[key1 key2]"))
		})
	})

	var _ = Describe("put", func() {
		var capturedBucket string
		var capturedKey string
		var capturedVal []byte

		makeRequest := func() *httptest.ResponseRecorder {
			api := kv.NewApi(lagertest.NewTestLogger("api"), &kv.StubStore{
				SetImpl: func(bucket string, key string, val []byte) error {
					capturedBucket = bucket
					capturedKey = key
					capturedVal = val
					return nil
				},
			})

			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodPut, path+"/id1/key1", bytes.NewBufferString("some value"))

			api.ServeHTTP(recorder, request)

			return recorder
		}

		It("returns 204", func() {
			response := makeRequest()

			Expect(response.Code).To(Equal(http.StatusNoContent))
		})

		It("sends to store", func() {
			Expect(capturedBucket).To(Equal("id1"))
			Expect(capturedKey).To(Equal("key1"))
			Expect(string(capturedVal)).To(Equal("some value"))
		})
	})

	var _ = Describe("get", func() {
		const presentKey = "key1"
		const missingKey = "nope"

		makeRequest := func(key string) *httptest.ResponseRecorder {
			api := kv.NewApi(lagertest.NewTestLogger("api"), &kv.StubStore{
				GetImpl: func(string, key string) ([]byte, error) {
					var val []byte
					if key == presentKey {
						val = []byte("value1")
					}
					return val, nil
				},
			})

			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, path+"/id1/"+key, nil)

			api.ServeHTTP(recorder, request)

			return recorder
		}

		It("sets content type", func() {
			response := makeRequest(presentKey)
			contentType := response.Header().Get("Content-Type")
			Expect(contentType).To(Equal("text/plain"))
		})

		It("responds with 200", func() {
			response := makeRequest(presentKey)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("responds with {key:value} object", func() {
			response := makeRequest(presentKey)

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("value1"))
		})

		It("responds with 404 if not found", func() {
			response := makeRequest(missingKey)

			Expect(response.Code).To(Equal(http.StatusNotFound))
		})
	})

	var _ = Describe("delete", func() {
		var capturedBucket string
		var capturedKey string

		makeRequest := func() *httptest.ResponseRecorder {
			api := kv.NewApi(lagertest.NewTestLogger("api"), &kv.StubStore{
				DelImpl: func(bucket string, key string) error {
					capturedBucket = bucket
					capturedKey = key
					return nil
				},
			})

			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodDelete, path+"/id1/key1", bytes.NewBufferString("some value"))

			api.ServeHTTP(recorder, request)

			return recorder
		}

		It("returns 204", func() {
			response := makeRequest()

			Expect(response.Code).To(Equal(http.StatusNoContent))
		})

		It("deletes from store", func() {
			Expect(capturedBucket).To(Equal("id1"))
			Expect(capturedKey).To(Equal("key1"))
		})
	})
})
