package main

import (
	"event_trigger/routes"
)

var (
	apiVersion = "v1"
)

func main() {
	a := routes.New(apiVersion)

	a.Start(a.Engine)
}
