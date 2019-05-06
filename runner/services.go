package main

import (
	"crypto/md5"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spiral/roadrunner/service"
	"github.com/spiral/roadrunner/service/rpc"
	"strings"
	"time"
)

//CustomServiceID _
const CustomServiceID = "custom"
const MaxFails = 10

// CustomConfig for service
type KafkaConfig struct {
	Brokers string
	Topic   string
	Ack     uint
}

// Hydrate config instance from .rr.* content
func (c *KafkaConfig) Hydrate(cfg service.Config) error {
	return cfg.Unmarshal(&c)
}

// CustomService _
type KafkaService struct {
	topic    string
	brokers  []string
	producer sarama.AsyncProducer
	errors []error
	Logger   *logrus.Logger
}

// Init service
func (s *KafkaService) Init(r *rpc.Service, cfg *KafkaConfig) (ok bool, err error) {

	sarama.Logger = s.Logger

	s.brokers = strings.Split(cfg.Brokers, ",")
	s.Logger.Printf("Kafka brokers: %s", strings.Join(s.brokers, ", "))
	s.topic = cfg.Topic

	err = r.Register("kafka", s)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Serve to start kafka service
func (s *KafkaService) Serve() error {
	config := sarama.NewConfig()
	config.Net.KeepAlive = 10 * time.Second
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	producer, err := sarama.NewAsyncProducer(s.brokers, config)
	if err != nil {
		return errors.Wrap(err, "failed to start Sarama producer")
	}
	s.producer = producer
	go func() {
		for range producer.Successes() {
			s.errors = []error{}
		}
	}()
	go func() {
		for err := range producer.Errors() {
			s.errors = append(s.errors, err)
			s.Logger.Warnln("Failed to write access log entry:", err)
		}
	}()
	return nil
}

// Stop kafka service
func (s *KafkaService) Stop() {
	err := s.producer.Close()
	if err != nil {
		s.Logger.Warnln(err)
	}
	return
}

func (s *KafkaService) Produce(message string, output *string) error {
	if len(s.errors) >= MaxFails {
		err := fmt.Errorf("kafka delivery errors %s > %s: %s", len(s.errors), MaxFails, s.errors[0])
		s.Logger.Warnln(err)
		return err
	}
	h := md5.New()
	s.producer.Input() <- &sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.ByteEncoder(h.Sum([]byte(message))),
		Value: sarama.StringEncoder(message),
	}
	return nil
}
