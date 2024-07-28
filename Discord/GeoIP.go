package Discord

import (
	"encoding/json"
	"io"
	"net/http"
)

type GeoIP struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Region      string `json:"region"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	ASN         string `json:"asn"`
}

func GetGeoIP(ip string) (geoIp GeoIP, err error) {
	req, err := http.NewRequest(http.MethodGet, "https://ipwhois.app/json/"+ip, nil)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &geoIp)
	if err != nil {
		return
	}
	return
}
