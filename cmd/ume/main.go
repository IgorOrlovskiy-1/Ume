package main

import (
	"Ume/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	//TODO: logger

	//TODO: bd

	//TODO: router

	//TODO: middlewars

	//TODO: tests
}
