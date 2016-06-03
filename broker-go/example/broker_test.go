package example_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cholick/cf-examples/broker-go/example"
	"github.com/cholick/cf-examples/broker-go/kv"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/brokerapi/Godeps/_workspace/src/github.com/pborman/uuid"
)

var _ = Describe("Broker", func() {
	const instanceId = "90dc5862-8e35-4046-b3b9-860b60174a9c"
	const orgId = "89364994-6b78-4dfd-8b7f-307c40330513"
	const spaceId = "28992410-3895-4cd7-b1a2-9c304f5aab4c"

	var broker *example.SampleServiceBroker
	var provisionDetails brokerapi.ProvisionDetails

	BeforeEach(func() {
		provisionDetails = brokerapi.ProvisionDetails{
			ServiceID:        example.ServiceId,
			PlanID:           example.BasicPlanId,
			OrganizationGUID: orgId,
			SpaceGUID:        spaceId,
			RawParameters:    nil,
		}
	})

	It("catalog has service & plan", func() {
		broker = example.NewSampleServiceBroker(&kv.StubStore{})

		services := broker.Services()

		Expect(len(services)).To(Equal(1))
		Expect(len(services[0].Plans)).To(Equal(1))
	})

	Describe("provision", func() {
		It("provision creates bucket", func() {
			var capturedBucketName string
			broker = example.NewSampleServiceBroker(&kv.StubStore{
				CreateBucketImpl: func(bucketName string) {
					capturedBucketName = bucketName
				},
				SetImpl: func(string, string, []byte) error {
					return nil
				},
			})

			broker.Provision(instanceId, provisionDetails, true)

			Expect(capturedBucketName).To(Equal(instanceId))
		})

		It("generates & stores credentials", func() {
			var capturedBucket string
			var capturedKey string
			var capturedVal []byte

			broker = example.NewSampleServiceBroker(&kv.StubStore{
				SetImpl: func(bucket string, key string, val []byte) error {
					capturedBucket = bucket
					capturedKey = key
					capturedVal = val
					return nil
				},
			})

			broker.Provision(instanceId, provisionDetails, true)

			Expect(capturedBucket).To(Equal(example.CredentialsBucket))
			Expect(capturedKey).To(Equal(instanceId))
			Expect(capturedVal).NotTo(BeNil())
		})
	})

	Describe("provision", func() {
		It("returns credentials", func() {
			var capturedBucketName string
			broker = example.NewSampleServiceBroker(&kv.StubStore{
				CreateBucketImpl: func(bucketName string) {
					capturedBucketName = bucketName
				},
				GetImpl: func(bucket string, key string) ([]byte, error) {
					if bucket == example.CredentialsBucket && key == instanceId {
						return []byte("correct-horse-battery-staple"), nil
					}
					panic("Unexpected bucket/key")
				},
			})

			binding, err := broker.Bind(instanceId, uuid.New(), brokerapi.BindDetails{})

			Expect(err).To(BeNil())

			credentials := binding.Credentials.(example.SampleServiceCredentials)

			Expect(credentials.ID).To(Equal(instanceId))
			Expect(credentials.Password).To(Equal("correct-horse-battery-staple"))
		})
	})
})
