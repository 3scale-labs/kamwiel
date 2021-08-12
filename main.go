package main

import (
	"github.com/3scale-labs/kamwiel/pkg/interfaces/controllers"
	"github.com/3scale-labs/kamwiel/pkg/interfaces/http"
)

func main() {
	go http.Start()
	controllers.Start()
}
