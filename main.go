package main

import (
	"os"
	"os/signal"
)

func main() {
	sm := NewServiceManager()
	sm.AddServices(
		NewTxnWorker(0, "http://192.168.0.211:4822", 10),
		NewGRPCServer(),
		NewRestfulServer(),
	)

	sm.StartServices()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	sm.StopServices()
}
