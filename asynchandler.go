package main

type Response struct {
	Query        string   `json:query`
	Results map[string]SingleResult `json:results`
}

type SingleResult struct {
	Url  string `json:url`
	Text string `json:text`
}
func interfaceFn(qp string, urls Url, tokens Tokens, ch chan map[string]SingleResult,
f func(string, Url, Tokens, chan map[string]SingleResult)) {
	f(qp, urls, tokens, ch)
}
func asyncHttpGets(qp string, config Configuration) (Response) {
	responses := map[string]SingleResult{}

	m := map[string]func(string, Url, Tokens, chan map[string]SingleResult){
		"google": google,
		"twitter": twitter,
		"duckduckgo": duckduckgo,
	}
	ch := make(chan map[string]SingleResult)
	for _, element := range config.Engines {
		go interfaceFn(qp, config.Urls, config.Tokens, ch, m[element])
	}
	for {
	      select {
	      case r := <-ch:
		      for k, v := range r {
			      responses[k] = v
		      }
		      if len(responses) == len(config.Engines) {
			      return Response{Query: qp, Results:responses}
		      }
	      }
	}
	//return Response{}
	return Response{Query: qp, Results:responses}
}
