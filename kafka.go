package main

import "github.com/spiral/roadrunner/service/rpc"

const ID = "roadrunner_symfony_kafka_example"

type Service struct{
}

func(s *Service) Init(r *rpc.Service) (ok bool, err error) {
	err 	= r.Register("custom", &rpcService{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func(s *Service) Serve() error {
	return nil
}

func(s *Service) Stop() {
	return
}

type rpcService struct {

}

func (s *rpcService) Hello(input string, output *string) error {
	*output = input
	return nil
}