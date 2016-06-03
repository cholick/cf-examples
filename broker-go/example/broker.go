package example

import (
	"errors"
	"github.com/cholick/cf-examples/broker-go/kv"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/brokerapi/Godeps/_workspace/src/github.com/pborman/uuid"
)

type SampleServiceBroker struct {
	store kv.Store
}

type SampleServiceCredentials struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	//todo: api url?
}

const ServiceId = "c44a4ede-1695-4194-897c-0197fbffa0f5"
const BasicPlanId = "087e0dca-2899-49ae-8ff5-46461de8f4bd"
const CredentialsBucket = "credentials"

func NewSampleServiceBroker(store kv.Store) *SampleServiceBroker {
	return &SampleServiceBroker{
		store: store,
	}
}

func (*SampleServiceBroker) Services() []brokerapi.Service {
	free := true
	plan := brokerapi.ServicePlan{
		ID:          BasicPlanId,
		Name:        "basic",
		Description: "Basic free plan",
		Free:        &free,
		Metadata:    nil,
	}
	service := brokerapi.Service{
		ID:            ServiceId,
		Name:          "kv_store",
		Description:   "A simple http key-value store",
		Bindable:      true,
		Tags:          []string{"key-value", "persistence"},
		PlanUpdatable: false,
		Plans: []brokerapi.ServicePlan{
			plan,
		},
	}
	return []brokerapi.Service{service}
}

func (this *SampleServiceBroker) Provision(
	instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool,
) (
	brokerapi.ProvisionedServiceSpec, error,
) {
	this.store.CreateBucket(instanceID)
	this.store.Set(CredentialsBucket, instanceID, []byte(uuid.New()))

	return brokerapi.ProvisionedServiceSpec{
		IsAsync:      false,
		DashboardURL: "",
	}, nil
}

func (this *SampleServiceBroker) Deprovision(
	instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (
	brokerapi.IsAsync, error,
) {
	//todo: delete bucket & credentials
	return false, nil
}

func (this *SampleServiceBroker) LastOperation(instanceID string) (brokerapi.LastOperation, error) {
	//no-op
	return brokerapi.LastOperation{}, nil
}

func (this *SampleServiceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	val, _ := this.store.Get(CredentialsBucket, instanceID)

	return brokerapi.Binding{
		Credentials: SampleServiceCredentials{
			ID:       instanceID,
			Password: string(val),
		},
		SyslogDrainURL:  "",
		RouteServiceURL: "",
	}, nil
}

func (*SampleServiceBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	//no-op
	return nil
}

func (this *SampleServiceBroker) Update(
	instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (
	brokerapi.IsAsync, error,
) {
	return false, errors.New("Unsupported operation")
}
