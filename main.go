package main

import (
	"os"
	"os/signal"
)

const priurl = "http://192.168.0.211:4822"
const puburl = "https://mainnet.infura.io/v3/9fb53ab19227473db75b4aca7c34cf3f"

func main() {
	sm := NewServiceManager()
	sm.AddServices(
		NewTxnWorker(0, priurl, 10),
		NewGRPCServer(),
		NewRestfulServer(),
	)

	sm.StartServices()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	sm.StopServices()
}
