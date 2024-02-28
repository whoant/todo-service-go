package rpc_caller

import (
	"flag"

	"todo-service/go-sdk/logger"
)

type apiItemCaller struct {
	name       string
	serviceUrl string
	logger     logger.Logger
}

func NewApiItemCaller(name string) *apiItemCaller {
	return &apiItemCaller{
		name: name,
	}
}

func (caller *apiItemCaller) GetPrefix() string {
	return caller.name
}

func (caller *apiItemCaller) Get() interface{} {
	return caller
}

func (caller *apiItemCaller) Name() string {
	return caller.name
}

func (caller *apiItemCaller) InitFlags() {
	flag.StringVar(&caller.serviceUrl, "item-service-url", "http://localhost:3001", "Url of item service")
}

func (caller *apiItemCaller) Configure() error {
	caller.logger = logger.GetCurrent().GetLogger("api.item")

	return nil
}

func (caller *apiItemCaller) Run() error {
	return nil
}

func (caller *apiItemCaller) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

func (caller *apiItemCaller) GetServiceUrl() string {
	return caller.serviceUrl
}
