package main

import (
	"github.com/3scale-labs/kamwiel/pkg/controllers"
	"github.com/3scale-labs/kamwiel/pkg/http"
)

func main() {
	go http.Start()
	controllers.Start()
}
