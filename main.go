package main

import (
	"TaxenHeimer/Discord"
	"TaxenHeimer/Minecraft"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tebeka/atexit"
	"log"
	"math/rand"
	"sync"
	"time"
)

func FatalD(error error) {
	fmt.Println(error.Error())
	Discord.NotifyError(error.Error())
}

func handleSignals() {
	err := Discord.Notification("Server Scanner has stopped.", 0xFF0000, "TaxenHeimer Status")
	if err != nil {
		FatalD(err)
	}
}

func randomIP() uint32 {
	var ip uint32 = 0
	for i := 0; i < 4; i++ {
		ip |= uint32(rand.Intn(255)) << (uint32(i) * 8)
	}
	return ip
}

func intToIP(ip uint32) string {
	firstOctet := (ip >> 24) & 0xFF
	secondOctet := (ip >> 16) & 0xFF
	thirdOctet := (ip >> 8) & 0xFF
	fourthOctet := ip & 0xFF
	ipStr := fmt.Sprintf("%d.%d.%d.%d", firstOctet, secondOctet, thirdOctet, fourthOctet)
	return ipStr
}

func lookupServer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ip := intToIP(randomIP())
			var port uint16 = 25565
			server, err := Minecraft.NewMinecraftServer(ip, port, "1.20.2")
			if err != nil {
				continue
			}
			info, err := server.GetStatusRequest()
			if err != nil {
				continue
			}
			err = Discord.NotifyServer(server.GetIP(), server.GetPort(), info)
			if err != nil {
				FatalD(err)
				return
			}
		}
	}
}

func main() {
	THREADS := 30
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		return
	}
	atexit.Register(handleSignals)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go lookupServer(ctx)
	}

	wg.Wait()
}
