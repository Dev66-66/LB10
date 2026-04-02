package main

import (
	"log"

	"github.com/Dev66-66/LB10/go-service/internal/app"
)

func main() {
	a := app.New()
	log.Println("msg=starting server addr=:8080")
	if err := a.Run(":8080"); err != nil {
		log.Fatalf("msg=server failed err=%v", err)
	}
}
