package main

import (
	"encoding/json"
	//"log"
	"net/http"
	"fmt"
	//"os"
	"io/ioutil"
)
type DDGResponse struct {
	RelatedTopics []RelatedTopic `json:RelatedTopics`
}
type RelatedTopic struct{
	FirstURL string `json:FirstURL`
	Text string `json:Text`
}
func duckduckgo(qp string, uri Url, tokens Tokens, ch chan map[string]SingleResult){
	url := uri.Duckduckgo
	client := &http.Client{}
	var ddgResult SingleResult
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ddgResult = SingleResult{"404", err.Error()}
	}

	q := req.URL.Query()
	q.Add("q", qp)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		ddgResult = SingleResult{"404", err.Error()}
	}
	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	var ddgResponse DDGResponse
	err = json.Unmarshal(resp_body, &ddgResponse)

	if err != nil {
		ddgResult = SingleResult{"404", err.Error()}
	}
	if len(ddgResponse.RelatedTopics) > 0 {
		ddgResult = SingleResult{ddgResponse.RelatedTopics[0].FirstURL, ddgResponse.RelatedTopics[0].Text}
	}else{
		ddgResult = SingleResult{"400", "Empty"}
	}
	ch <- map[string]SingleResult{"duckduckgo": ddgResult}
}