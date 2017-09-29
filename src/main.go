package main

import (
	"model"
)

const (
	SOCKET = ":8000"
)

func main() {
	a := &models.App{Socket: SOCKET}
	a.Run()
}
