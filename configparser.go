package main

import (
    "fmt"
    "encoding/json"
    "os"
)

type Configuration struct {
    Engines []string `json:engines`
    Urls Url `json:urls`
    Tokens Tokens `json:tokens`
}
type Url struct{
	Google string `json:google`
	Twitter string `json:twitter`
	Duckduckgo string `json:duckduckgo`
}
type Tokens struct{
	Google Google `json:google`
	Twitter Twitter `json:twitter`
}
type Google struct{
	Key string `json:key`
	Cx string `json:cx`
}
type Twitter struct{
	Consumer_key string `json:consumer_key`
	Consumer_secret string `json:consumer_secret`
	Access_token string `json:access_token`
	Access_token_secret string `json:access_token_secret`
}
func getConfig() Configuration{
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
	  fmt.Println("error:", err)
	}
	return configuration
}