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
	cep := "66825-070"

	go consultarCEPViaCep(channelViaCep, &c, cep)
	go consultarCEPCdn(channelCDN, &c, cep)

	result := getResult(channelViaCep, channelCDN)

	fmt.Println(result)
}

func consultarCEPViaCep(ch chan<- string, c *http.Client, cep string) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	resp, err := c.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response, _ := io.ReadAll(resp.Body)
	ch <- string(response)
}

func consultarCEPCdn(ch chan<- string, c *http.Client, cep string) {
	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
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
