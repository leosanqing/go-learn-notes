package main

import "go.uber.org/zap"

func main() {
	log,_ :=zap.NewProduction()
	log.Warn("this is a test")
}
