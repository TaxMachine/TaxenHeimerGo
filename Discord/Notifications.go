package Discord

import (
	"TaxenHeimer/Minecraft"
	"fmt"
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
	embed := NewEmbed(
		"```"+net.JoinHostPort(ip, strconv.Itoa(int(port)))+"```",
		"```"+info.Description.Text+"```",
		0x00FF00)
	embed.AddField("Version", info.Version.Name, true)
	embed.AddField("Software", info.Version.Software, true)
	embed.AddField("Protocol", string(info.Version.Protocol), true)
	embed.AddField("Players", string(info.Players.Online)+"/"+string(info.Players.Max), true)
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
	req := WebhookRequest{
		Content:  "caca",
		Username: "TaxenHeimer",
		TTS:      false,
		Embeds:   []Embed{embed},
	}
	res, err := SendWebhook(os.Getenv("serverwebhook"), req)
	if !res {
		return fmt.Errorf("failed to send to webhook")
	}
	return nil
}
