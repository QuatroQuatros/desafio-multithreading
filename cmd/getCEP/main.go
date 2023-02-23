package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cep := os.Args[1:]

	if len(cep) < 1 {
		panic("You must pass a zip code with a parameter")
	}
	buscaCEP(cep[0])
}

func buscaCEP(cep string) {
	ch1 := make(chan []byte)
	ch2 := make(chan []byte)

	go func() {
		//time.Sleep(time.Second * 1)
		req, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
		if err != nil {
			log.Fatal("Erro ao fazer requisição")
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println("Erro ao ler resposta")
		}
		ch1 <- res
	}()

	go func() {
		// time.Sleep(time.Second * 1)
		req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

		if err != nil {
			log.Fatal("Erro ao fazer requisição")
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println("Erro ao ler resposta")
		}
		ch2 <- res
	}()

	select {
	case msg := <-ch1:
		fmt.Printf("Resposta recebida do endereço https://cdn.apicep.com/file/apicep/%s.json\n %s\n", cep, msg)

	case msg := <-ch2:
		fmt.Printf("Resposta recebida do endereço http://viacep.com.br/ws/%s/json/\n %s\n", cep, msg)

	case <-time.After(time.Second * 1):
		fmt.Println("Request timeout")
	}

}
