package Discord

import (
	"TaxenHeimer/Minecraft"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func NotifyServer(ip string, port uint16, info Minecraft.ServerInfo) (err error) {
	geoIp, err := GetGeoIP(ip)
	if err != nil {
		return
	}
	embed := NewEmbed()
	embed.SetTitle("```" + net.JoinHostPort(ip, strconv.Itoa(int(port))) + "```")
	embed.SetDescription("```" + info.Description.Text + "```")
	embed.SetColor(0x00FF00)
	embed.AddField("Version", info.Version.Name, true)
	embed.AddField("Software", info.Version.Software, true)
	embed.AddField("Protocol", strconv.Itoa(int(info.Version.Protocol)), true)
	embed.AddField("Players", strconv.Itoa(int(info.Players.Online))+"/"+strconv.Itoa(int(info.Players.Max)), true)
	embed.AddField("Country", geoIp.Country+" :flag_"+strings.ToLower(geoIp.CountryCode)+":", true)
	embed.AddField("Region", geoIp.Region, true)
	embed.AddField("City", geoIp.City, true)
	embed.AddField("ISP", geoIp.ISP, true)
	embed.AddField("ASN", geoIp.ASN, true)
	if len(info.Players.Sample) != 0 {
		players := "```\n"
		for _, pe := range info.Players.Sample {
			players += pe.Name + " : " + pe.Id + "\n"
		}
		players += "```"
		embed.AddField("Players", players, true)
	}

	serverWebhook, err := NewWebhook(os.Getenv("serverwebhook"))
	if err != nil {
		return
	}
	serverWebhook.SetUsername("TaxenHeimer")
	serverWebhook.SetMessage("caca")
	serverWebhook.AddEmbed(embed)
	err = serverWebhook.Send()
	if err != nil {
		return fmt.Errorf("failed to send to webhook: %s", err)
	}
	return
}

func Notification(message string, color int, username string) (err error) {
	embed := NewEmbed()
	embed.SetTitle(message)
	embed.SetColor(color)

	statusWebhook, err := NewWebhook(os.Getenv("statuswebhook"))
	if err != nil {
		return
	}
	statusWebhook.SetUsername(username)
	statusWebhook.AddEmbed(embed)
	err = statusWebhook.Send()
	if err != nil {
		return fmt.Errorf("failed to send to webhook: %s", err)
	}
	return
}

func NotifyError(message string) {
	errorwebhook, err := NewWebhook(os.Getenv("errorwebhook"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	errorwebhook.SetMessage("```\n" + message + "\n```")
	err = errorwebhook.Send()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
