package main

import (
	"os"
	"os/signal"
)

func main() {
	sm := NewServiceManager()
	sm.AddServices(
		NewTxnWorker(0, "https://mainnet.infura.io/v3/9fb53ab19227473db75b4aca7c34cf3f", 10),
		NewGRPCServer(),
		NewRestfulServer(),
	)

	sm.StartServices()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	sm.StopServices()
}
