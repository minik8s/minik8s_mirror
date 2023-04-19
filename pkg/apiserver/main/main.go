package main

import (
	"miniK8s/pkg/apiserver/app"
	"miniK8s/pkg/apiserver/config"
)

func main() {
	apiServer := apiserver.New(config.DefaultServerConfig())
	apiServer.Run()
}