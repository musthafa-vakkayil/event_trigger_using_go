package main

import (
	_ "event_trigger/docs"
	"event_trigger/routes"
)

var (
	apiVersion = "v1"
)

// @title Event Trigger API
// @version 1.0
// @description This is an API for the Event Trigger application.

// @host 18.205.239.53:8080
// @BasePath /v1
func main() {

	a := routes.New(apiVersion)

	a.Start(a.Engine)
}
