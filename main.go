package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CEPResponse struct {
	API string
	CEP string
}

func main() {
	c := http.Client{}
	channelViaCep := make(chan string)
	channelCDN := make(chan string)

	go consultarCEP("https://viacep.com.br/ws/66825-070/json/", channelViaCep, &c)
	go consultarCEP("https://cdn.apicep.com/file/apicep/66825-070.json", channelCDN, &c)

	result := getResult(channelViaCep, channelCDN)

	fmt.Println(result)
}

func consultarCEP(url string, ch chan<- string, c *http.Client) {
	resp, err := c.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response, _ := io.ReadAll(resp.Body)
	ch <- string(response)
}

func getResult(channelViaCep chan string, channelCDN chan string) CEPResponse {
	select {
	case msgViaCep := <-channelViaCep:
		return CEPResponse{API: "http://viacep.com.br/ws/", CEP: msgViaCep}
	case msgCDN := <-channelCDN:
		return CEPResponse{API: "https://cdn.apicep.com/file/apicep/", CEP: msgCDN}
	case <-time.After(time.Second * 5):
		panic(errors.New("endpoint request timeout"))
	}
}
