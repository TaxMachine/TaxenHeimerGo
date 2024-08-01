package main

import (
	"TaxenHeimer/Discord"
	"TaxenHeimer/Minecraft"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tebeka/atexit"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	OvhIpRanges = []string{
		"107.189.64.0", "135.125.0.0", "135.125.128.0", "135.148.0.0", "135.148.128.0", "137.74.0.0", "139.99.0.0", "139.99.128.0", "141.227.128.0",
		"141.227.130.0", "141.227.132.0", "141.227.134.0", "141.227.136.0", "141.227.138.0", "141.227.140.0", "141.94.0.0", "141.95.0.0", "141.95.128.0",
		"142.4.192.0", "142.44.128.0", "142.44.140.0", "144.217.0.0", "145.239.0.0", "146.59.0.0", "146.59.0.0", "147.135.0.0", "147.135.128.0", "148.113.0.0",
		"148.113.128.0", "149.202.0.0", "149.56.0.0", "151.80.0.0", "15.204.0.0", "15.204.128.0", "152.228.128.0", "15.235.0.0", "15.235.128.0", "158.69.0.0",
		"162.19.0.0", "162.19.128.0", "164.132.0.0", "167.114.0.0", "167.114.128.0", "167.114.192.0", "176.31.0.0", "178.32.0.0", "185.45.160.0", "188.165.0.0",
		"192.240.152.0", "192.31.246.0", "192.95.0.0", "192.99.0.0", "192.99.65.0", "193.70.0.0",
	}
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

func randomIP() int {
	var ip = 0
	for i := 0; i < 4; i++ {
		ip |= rand.Intn(255) << (i * 8)
	}
	return ip
}

func intToIP(ip int) string {
	firstOctet := (ip >> 24) & 0xFF
	secondOctet := (ip >> 16) & 0xFF
	thirdOctet := (ip >> 8) & 0xFF
	fourthOctet := ip & 0xFF
	ipStr := fmt.Sprintf("%d.%d.%d.%d", firstOctet, secondOctet, thirdOctet, fourthOctet)
	return ipStr
}

func IPToInt(ip string) uint32 {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0
	}

	var ipInt uint32
	for i, part := range parts {
		octet, err := strconv.Atoi(part)
		if err != nil || octet < 0 || octet > 255 {
			return 0
		}
		ipInt |= uint32(octet) << uint32(24-i*8)
	}

	return ipInt
}

func enumerateServer(startIp int) {
	for {
		if startIp&0xFF == 255 {
			return
		}
		fmt.Printf("Trying %s\n", intToIP(startIp))
		server, err := Minecraft.NewMinecraftServer(intToIP(startIp), 25565, "1.20.2")
		if err != nil {
			startIp++
			continue
		}
		info, err := server.GetStatusRequest()
		if err != nil {
			startIp++
			continue
		}
		server.Close()
		fmt.Println("VALID", intToIP(startIp))
		err = Discord.NotifyServer(server.GetIP(), server.GetPort(), info)
		if err != nil {
			FatalD(err)
			continue
		}
		startIp++
		time.Sleep(2 * time.Second)
	}
}

func main() {
	//MaxAddress := math.Pow(2, 32)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		return
	}
	atexit.Register(handleSignals)

	for _, ip := range OvhIpRanges {
		enumerateServer(int(IPToInt(ip)))
	}
}
