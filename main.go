package main

import (
	"TaxenHeimer/Discord"
	"TaxenHeimer/Minecraft"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
	ip := "join.3b3france.fr"
	var port uint16 = 25011
	server, err := Minecraft.NewMinecraftServer(ip, port, "1.20.2")
	if err != nil {
		fmt.Println("Connection: ", err)
		return
	}
	info, err := server.GetStatusRequest()
	if err == nil {
		err = Discord.NotifyServer(ip, port, info)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
