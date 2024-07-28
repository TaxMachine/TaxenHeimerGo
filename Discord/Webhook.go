package Discord

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WebhookResponse struct {
	Name          string `json:"name"`
	Id            string `json:"id"`
	ChannelId     string `json:"channel_id"`
	GuildId       string `json:"guild_id"`
	ApplicationId string `json:"application_id"`
	Token         string `json:"token"`
	Url           string `json:"url"`
	Avatar        string `json:"avatar"`
	Type          int    `json:"type"`
}

type WebhookRequest struct {
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarUrl string  `json:"avatar_url"`
	TTS       bool    `json:"tts"`
	Embeds    []Embed `json:"embeds"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Url         string       `json:"url"`
	Timestamp   string       `json:"timestamp"`
	Color       int          `json:"color"`
	Footer      EmbedFooter  `json:"footer"`
	Image       EmbedImage   `json:"image"`
	Author      EmbedAuthor  `json:"author"`
	Fields      []EmbedField `json:"fields"`
}

type EmbedFooter struct {
	Text    string `json:"text"`
	IconUrl string `json:"icon_url"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedAuthor struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"icon_url"`
}

type EmbedImage struct {
	Url string `json:"url"`
}

func NewEmbed(title string, description string, color int) (embed Embed) {
	embed = Embed{
		Title:       title,
		Description: description,
		Color:       color,
	}
	return
}

func (embed *Embed) AddField(name string, value string, inline bool) {
	embed.Fields = append(embed.Fields, EmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}

func SendWebhook(url string, body WebhookRequest) (success bool, err error) {
	jsonBody, err := json.Marshal(&body)
	bodyBytes := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, url, bodyBytes)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	return res.StatusCode == 200 || res.StatusCode == 204, nil
}
