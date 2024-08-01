package Discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

type WebhookError struct {
	Message string `json:"message"`
	Code    int32  `json:"code"`
}

type WebhookRequest struct {
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarUrl string  `json:"avatar_url"`
	TTS       bool    `json:"tts"`
	Embeds    []Embed `json:"embeds"`
}

type Webhook struct {
	Url  string
	Body WebhookRequest
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

func NewEmbed() (embed Embed) {
	embed = Embed{}
	return
}

func (embed *Embed) SetTitle(title string) {
	embed.Title = title
}

func (embed *Embed) SetDescription(description string) {
	embed.Description = description
}

func (embed *Embed) SetColor(color int) {
	embed.Color = color
}

func (embed *Embed) AddField(name string, value string, inline bool) {
	if len(embed.Fields) == 25 {
		return
	}
	embed.Fields = append(embed.Fields, EmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}

func (embed *Embed) AddAuthor(name string, url string, iconUrl string) {
	embed.Author = EmbedAuthor{
		Name:    name,
		Url:     url,
		IconUrl: iconUrl,
	}
}

func (embed *Embed) AddFooter(text string, iconUrl string) {
	embed.Footer = EmbedFooter{
		Text:    text,
		IconUrl: iconUrl,
	}
}

func NewWebhook(url string) (webhook Webhook, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return Webhook{}, fmt.Errorf("invalid webhook")
	}

	webhook.Url = url

	return
}

func (webhook *Webhook) SetUsername(username string) {
	webhook.Body.Username = username
}

func (webhook *Webhook) SetAvatar(avatarUrl string) {
	webhook.Body.AvatarUrl = avatarUrl
}

func (webhook *Webhook) SetMessage(message string) {
	webhook.Body.Content = message
}

func (webhook *Webhook) AddEmbed(embed Embed) {
	if len(webhook.Body.Embeds) == 10 {
		return
	}
	webhook.Body.Embeds = append(webhook.Body.Embeds, embed)
}

func (webhook *Webhook) Send() (err error) {
	jsonBody, err := json.Marshal(&webhook.Body)
	if err != nil {
		return
	}
	bodyBytes := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, webhook.Url, bodyBytes)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 204 {
		var werror WebhookError
		var body []byte
		body, err = io.ReadAll(res.Body)
		err = json.Unmarshal(body, &werror)
		if err != nil {
			return
		}
		if werror.Message == "You are being rate limited." {
			time.Sleep(2000)
			err = webhook.Send()
		}
		err = fmt.Errorf("webhook request error: %s", werror.Message)
	}
	return
}
