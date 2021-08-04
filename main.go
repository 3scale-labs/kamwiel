package main

import (
	"github.com/3scale-labs/kamwiel/controllers"
	"github.com/3scale-labs/kamwiel/pkg/routing"
)

func main() {
	go routing.Start()
	controllers.Start()
}
