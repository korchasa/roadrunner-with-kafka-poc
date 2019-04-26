package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spiral/roadrunner/service"
	"github.com/spiral/roadrunner/service/rpc"
	"strings"
)

//CustomServiceID _
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

// CustomService _
type CustomService struct {
	Logger      *logrus.Logger
	KafkaWriter *kafka.Writer
}

// Init service
func (s *CustomService) Init(r *rpc.Service, cfg *CustomConfig) (ok bool, err error) {
	brokerList := strings.Split(cfg.Brokers, ",")
	s.Logger.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	s.KafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokerList,
		Topic:   cfg.Topic,
		Async:   true,
	})

	err = r.Register("kafka", &kafkaService{
		writer: s.KafkaWriter,
		topic:  cfg.Topic,
	})

	if err != nil {
		logrus.Warnln(err)
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
	err := s.KafkaWriter.Close()
	if err != nil {
		logrus.Warnln(err)
	}
	return
}

type kafkaService struct {
	writer *kafka.Writer
	topic  string
}

func (s *kafkaService) Produce(message string, output *string) error {
	return s.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("0"),
		Value: []byte(message),
	})
}
