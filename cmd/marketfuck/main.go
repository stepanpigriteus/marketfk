package main

import (
	"fmt"
	"marketfuck/pkg/config"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)
}
