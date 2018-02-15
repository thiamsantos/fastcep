package cep

// CEP (postal address code)
type CEP struct {
	CEP          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

// CEPSize is the size of a cep
const CEPSize = 8
