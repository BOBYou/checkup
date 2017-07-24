package g

import (
	"strings"
	"net/http"
  "log"
  "io/ioutil"
  "net/url"
)


func Post(url string, urlData url.Values ) (string,error){

	s := urlData.Encode()

	req, err := http.NewRequest("POST", url, strings.NewReader(s))

	if err != nil {
		log.Printf("http.NewRequest() error: %v\n", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Printf("http.Do() error: %v\n", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll() error: %v\n", err)
	}

  return string(data),err
}
