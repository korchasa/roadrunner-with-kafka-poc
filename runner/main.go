package main

import (
	"github.com/sirupsen/logrus"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"

	// services (plugins)
	"github.com/spiral/roadrunner/service/env"
	"github.com/spiral/roadrunner/service/http"
	"github.com/spiral/roadrunner/service/rpc"
	"github.com/spiral/roadrunner/service/static"

	// additional commands and debug handlers
	_ "github.com/spiral/roadrunner/cmd/rr/http"
)

func main() {
	rr.Container.Register(env.ID, &env.Service{})
	rr.Container.Register(rpc.ID, &rpc.Service{})
	rr.Container.Register(http.ID, &http.Service{})
	rr.Container.Register(static.ID, &static.Service{})
	rr.Container.Register(CustomServiceID, &CustomService{
		Logger: rr.Logger,
	})

	rr.Logger.Formatter = &logrus.TextFormatter{ForceColors: true}

	rr.Execute()
}
