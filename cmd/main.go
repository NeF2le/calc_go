package main

import (
	"fmt"

	"github.com/NeF2le/calc_go/internal/application"
)

func main() {
	app := application.NewApplication()
	if err := app.RunServer(); err != nil {
		fmt.Println("Server terminated with error:", err)
	}
}
