package address

// Address struct
type Address struct {
	CEP          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Uf           string `json:"uf"`
}

// CEPSize is the size of a cep
const CEPSize = 8
