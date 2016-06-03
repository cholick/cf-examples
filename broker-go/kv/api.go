package kv

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pivotal-golang/lager"
)

func NewApi(logger lager.Logger, store Store) *mux.Router {
	router := mux.NewRouter()
	api := &KvApi{
		store: store,
	}

	//todo: credentials and auth

	router.HandleFunc("/kv/{id}/", api.List).Methods(http.MethodGet)
	router.HandleFunc("/kv/{id}/{key}", api.Set).Methods(http.MethodPut)
	router.HandleFunc("/kv/{id}/{key}", api.Get).Methods(http.MethodGet)
	router.HandleFunc("/kv/{id}/{key}", api.Del).Methods(http.MethodDelete)

	return router
}

type KvApi struct {
	store Store
}

func (this *KvApi) List(resp http.ResponseWriter, req *http.Request) {
	bucket, _ := parseVars(req)

	keys, _ := this.store.List(bucket)
	resp.Header().Set("Content-Type", "text/plain")
	resp.Write([]byte(fmt.Sprintf("%s", keys)))
}

func (this *KvApi) Set(resp http.ResponseWriter, req *http.Request) {
	bucket, key := parseVars(req)
	body, _ := ioutil.ReadAll(req.Body)

	this.store.Set(bucket, key, body)

	resp.WriteHeader(http.StatusNoContent)
}

func (this *KvApi) Get(resp http.ResponseWriter, req *http.Request) {
	bucket, key := parseVars(req)

	val, _ := this.store.Get(bucket, key)
	if val != nil {
		resp.Header().Set("Content-Type", "text/plain")
		resp.Write(val)
	} else {
		http.NotFound(resp, req)
	}
}

func (this *KvApi) Del(resp http.ResponseWriter, req *http.Request) {
	bucket, key := parseVars(req)

	this.store.Del(bucket, key)

	resp.WriteHeader(http.StatusNoContent)
}

func parseVars(req *http.Request) (bucket string, key string) {
	vars := mux.Vars(req)
	return vars["id"], vars["key"]
}
