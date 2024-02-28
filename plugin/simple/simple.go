package simple

import (
	"flag"
)

type simplePlugin struct {
	prefix string
	name   string
	value  string
}

func NewSimplePlugin(name, prefix string) *simplePlugin {
	return &simplePlugin{
		prefix: prefix,
		name:   name,
	}
}

func (s *simplePlugin) GetPrefix() string {
	return s.prefix
}

func (s *simplePlugin) Get() interface{} {
	return s
}

func (s *simplePlugin) Name() string {
	return s.name
}

func (s *simplePlugin) InitFlags() {
	prefix := s.prefix
	if prefix != "" {
		prefix += "-"
	}
	flag.StringVar(&s.value, prefix+"value", "default value", "Some value of simple plugin")
}

func (s *simplePlugin) Configure() error {
	return nil
}

func (s *simplePlugin) Run() error {
	return nil
}

func (s *simplePlugin) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()

	return c
}

func (s *simplePlugin) GetValue() string {
	return s.value
}
