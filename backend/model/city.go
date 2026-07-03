package model

// City representa uma capital lida do arquivo capitais.json.
//
// Atenção: o NOME da capital não fica aqui dentro — ele é a chave do objeto
// no JSON. Por isso, ao carregar o arquivo, usamos []map[string]City, onde a
// chave do mapa é o nome da capital. Ex.:
//
//	{ "Manaus": { "toll": 50, "neighbors": { "Boa Vista": 785 } } }
type City struct {
	// Toll é o pedágio cobrado ao passar por esta capital.
	Toll int `json:"toll"`
	// Neighbors mapeia capital_vizinha -> distância em km até ela.
	Neighbors map[string]int `json:"neighbors"`
}
