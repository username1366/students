package main

import (
	"models"
)

const (
	SOCKET = ":8000"
)

func main() {
	a := &models.App{Socket: SOCKET}
	a.Run()
}
