package main
import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
)
type GoogleResponse struct {
	Items []Item `json:items`
}
type Item struct{
	Link string `json:link`
	Snippet string `json:snippet`
}
func google(qp string, ur Url, token Tokens, ch chan map[string]SingleResult) {
	uri := ur.Google
	tokens := token.Google
	var googleResult SingleResult
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		googleResult = SingleResult{"404", err.Error()}
	}
	q := req.URL.Query()
	q.Add("q", qp)
	q.Add("key", tokens.Key)
	q.Add("cx", tokens.Cx)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		googleResult = SingleResult{"404", err.Error()}
	}
	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)
	var googleResponse *GoogleResponse
	err = json.Unmarshal(resp_body, &googleResponse)
	if err != nil {
		googleResult = SingleResult{"404", err.Error()}
	}
	if len(googleResponse.Items) > 0 {
		googleResult = SingleResult{googleResponse.Items[0].Link, googleResponse.Items[0].Snippet}
	}else{
		googleResult = SingleResult{"400", "Empty"}
	}
	ch <- map[string]SingleResult{"google": googleResult}
}
