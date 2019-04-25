package main

import "github.com/spiral/roadrunner/service/rpc"
import "github.com/spiral/roadrunner/service"
import "fmt"

const CustomServiceID = "custom"

// CustomConfig for service
type CustomConfig struct {
	Brokers string
	Topic   string
	Ack     uint
}

// Hydrate config instance from .rr.* content
func (c *CustomConfig) Hydrate(cfg service.Config) error {
	return cfg.Unmarshal(&c)
}

// CustomService
type CustomService struct {
}

// Init service
func (s *CustomService) Init(r *rpc.Service, cfg *CustomConfig) (ok bool, err error) {
	err = r.Register("kafka", &kafkaService{
		brokers: cfg.Brokers,
		topic:   cfg.Topic,
		ack:     cfg.Ack,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// Serve to start kafka service
func (s *CustomService) Serve() error {
	return nil
}

// Stop kafka service
func (s *CustomService) Stop() {
	return
}

type kafkaService struct {
	brokers string
	topic   string
	ack     uint
}

func (s *kafkaService) Produce(message string, output *string) error {
	*output = fmt.Sprintf("brokers: %s, topic: %s, ack: %d, message: %s", s.brokers, s.topic, s.ack, message)
	return nil
}
