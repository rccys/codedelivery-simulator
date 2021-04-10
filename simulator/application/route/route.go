package route

import ("bufio"
		"encoding/json"	
		"errors"
		"os"
		"strconv"
		"strings")

// definindo a estrutura de uma rota que iremos receber via kafka
type Route struct {
	ID string `json:"routeId"`// quando o ID da rota é 1 pegaremos o arquivo 1 da pasta destination
	ClientID string `json:"clientID"`
	Positions []Position `json:"position"` // e uma lista com diversas posicoes
}

type Position struct{
	Lat float64 `json:"lat"`
	Long float64 `json:"long"`
}

// struct com a forma do json que precisamos entregar para o Kafka
type PartialRoutePosition struct {
	// o go tem o recurso de tags que facilita a conversao para json
	ID string `json:"routeId"`
	ClientID string `json:"clientId"`
	Position []float64 `json:"position"`
	Finished bool `json:"finished"`// toda vez que for a ultima posicao o frontend sabera que a corrida encerrou
}

// vamos criar um metodo para carregar as posicoes para dentro da lista Position na rota 
func (r *Route) LoadPositions() error { 
	// no go o que vem antes da chave e o retorno, se nao ocorrer erro retorna vazio
	if r.ID == ""{
		return errors.New("rout id not informed")
	}
	// variavel com escopo pequeno pode ser abrviada, file vira f
	f, err := os.Open("destination/" + r.ID + ".txt")
	if err != nil { 
		return err
	}
	// espera tudo da funcao ser executada para fechar o arquivo
	defer f.Close()

	// le o conteudo de f
	scanner := bufio.NewScanner(f)
	// loop para pegar linha a linha
	for scanner.Scan(){
		data := strings.Split(scanner.Text(), ",") // split onde o separador e a virgula (separador nos arq 1,2 e 3)
		lat, err := strconv.ParseFloat(data[0], 64)
		if err != nil{
			return nil
		}
		long, err := strconv.ParseFloat(data[1], 64)
		if err != nil{
			return nil 
		}
		// agora vamos apensar na lista
		r.Positions = append(r.Positions, Position{
			Lat: lat,
			Long: long,
		})
	}
	return nil
}

// agora vamos criar uma funcao para gerar as posicoes em json e colocar numa lista de string
func (r *Route) ExportJsonPositions() ([]string, error) {
	var route PartialRoutePosition
	var result []string
	total := len(r.Positions)

	// vamos percorrer todas as posicoes dentro de []Poistion
	for k, v := range r.Positions{
		route.ID = r.ID
		route.ClientID = r.ClientID
		route.Position = []float64{v.Lat, v.Long}
		route.Finished = false
		if total-1 == k {
			route.Finished = true
		}
		// agora vamos converter para json
		jsonRoute, err := json.Marshal(route)
		if err != nil {
			return nil, err
		}
		// teremos uma lista de string com json em cada posicao para enviar pro frontend
		result = append(result, string(jsonRoute))
	}
	return result, nil
}