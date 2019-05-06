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
	"log"
	"os"
)

//CustomServiceID _
const CustomServiceID = "custom"
//MaxFails before service return error
const MaxFails = 100

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
	topic    string
	brokers  []string
	producer sarama.AsyncProducer
	errors []error
	successCount int
	Logger   *logrus.Logger
}

// Init service
func (s *CustomService) Init(r *rpc.Service, cfg *CustomConfig) (ok bool, err error) {

	// sarama.Logger = s.Logger

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
func (s *CustomService) Serve() error {
	config := sarama.NewConfig()
	config.ClientID = CustomServiceID
	config.Net.KeepAlive = 10 * time.Second
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(s.brokers, config)
	if err != nil {
		return errors.Wrap(err, "failed to start Sarama producer")
	}
	s.producer = producer

	lgr := log.New(os.Stdout, "[custom] ", log.LstdFlags)

	go func() {
		for range producer.Successes() {
			// s.errors = []error{}
			s.successCount++
		}
	}()
	go func() {
		for err := range producer.Errors() {
			s.errors = append(s.errors, err)
			lgr.Println("Logrus failed to write access log entry:", err)
			s.Logger.Println("Failed to write access log entry:", err)
		}
	}()
	return nil
}

// Stop kafka service
func (s *CustomService) Stop() {
	err := s.producer.Close()
	if err != nil {
		s.Logger.Warnln(err)
	}
	return
}

func (s *CustomService) Produce(message string, output *string) error {
	if len(s.errors) > MaxFails {
		err := fmt.Errorf("kafka delivery errors %d > %d: %s", len(s.errors), MaxFails, s.errors[0])
		s.Logger.Warnln(err)
		return err
	}
	h := md5.New()
	s.producer.Input() <- &sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.ByteEncoder(h.Sum([]byte(message))),
		Value: sarama.StringEncoder(message),
	}
	*output = fmt.Sprintf(
		"Failure:%d, Success:%d",
		len(s.errors),
		s.successCount,
	)
	return nil
}
