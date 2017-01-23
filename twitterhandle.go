package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
 )
type TwitterResponse struct{
	Statuses []Status `json:statuses`
}
type Status struct{
	Text string `json:text`
	User User `json:user`
}
type User struct{
	Url string `json:url`
}
func twitter(qp string, ur Url, token Tokens, ch chan map[string]SingleResult){
	uri := ur.Twitter
	tokens := token.Twitter
	var twitterResult SingleResult

	client := &http.Client{}
	encodedKeySecret := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		url.QueryEscape(tokens.Consumer_key),
		url.QueryEscape(tokens.Consumer_secret))))

	reqBody := bytes.NewBuffer([]byte(`grant_type=client_credentials`))
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", reqBody)
	if err != nil {
		log.Fatal(err, client, req)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedKeySecret))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, respBody)
	}

	type BearerToken struct {
		AccessToken string `json:"access_token"`
	}
	var b BearerToken
	json.Unmarshal(respBody, &b)
	req, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("q", qp)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", b.AccessToken))
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var twitterResponse TwitterResponse
	err = json.Unmarshal(respBody, &twitterResponse)
	if err != nil {
		fmt.Println("error:", err)
	}
	if len(twitterResponse.Statuses) > 0 {
		twitterResult = SingleResult{twitterResponse.Statuses[0].User.Url, twitterResponse.Statuses[0].Text}
	}else{
		twitterResult = SingleResult{"400", "Empty"}
	}
	ch <- map[string]SingleResult{"twitter": twitterResult}
}

