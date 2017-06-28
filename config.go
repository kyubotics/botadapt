package main

type Configuration struct {
	Host           string `json:"host"`
	Port           int `json:"port"`
	PostURL        string `json:"post_url"`
	MessageSources []map[string]string `json:"message_sources"`
}

var Config Configuration
