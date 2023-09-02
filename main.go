package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// https://mholt.github.io/json-to-go/ para converte json para stuc
type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	c1 := make(chan ViaCEP)
	c2 := make(chan BrasilApi)

	cep := strings.Join(os.Args[1:], " ")
	go coletaCep("http://www.viacep.com.br/ws/"+cep+"/json/", c1)
	go coletaCep("https://brasilapi.com.br/api/cep/v1/"+cep, c2)

	for {
		select {
		case msgViaCep := <-c1:
			msg := msgViaCep
			fmt.Println("Recebido de ViaCep")
			fmt.Printf("Cep: %s\nLogradouro: %s\nComplemento: %s\nBairro: %s\nLocalidade: %s\nUf: %s\nIbge: %s\nGia: %s\nDdd: %s\nSiafi: %s\n",
				msg.Cep,
				msg.Logradouro,
				msg.Complemento,
				msg.Bairro,
				msg.Localidade,
				msg.Uf,
				msg.Ibge,
				msg.Gia,
				msg.Ddd,
				msg.Siafi)
		case msgBrasilApi := <-c2:
			msg := msgBrasilApi
			fmt.Println("Recebido de BrasilApi")
			fmt.Printf("Cep: %s\nStreet: %s\nCity: %s\nNeighborhood: %s\nState: %s\nService: %s\n",
				msg.Cep,
				msg.Street,
				msg.City,
				msg.Neighborhood,
				msg.State,
				msg.Service)
		case <-time.After(time.Second * 1):
			println("timeout")
		}
		break
	}
}

func coletaCep[t ViaCEP | BrasilApi](url string, c chan t) {
	req, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição em %v", url)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta de %v", url)
	}
	var data t
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta de %v", url)
	}
	c <- data
}
