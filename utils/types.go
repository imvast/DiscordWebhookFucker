package utils

type Config struct {
	Threads int `json:"threads"`
	Debug bool `json:"debug"`
	Webhook string `json:"webhook"`
    Content string `json:"content"`
	AvatarUrl string `json:"avatar_url"`
}