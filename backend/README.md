# Caminho Mais Barato entre Capitais — Backend

Backend em Go para o trabalho de **Teoria de Grafos (UNESC)**.

O sistema modela um **grafo não direcionado** das capitais brasileiras a partir do arquivo `api/capitais.json` e disponibiliza uma API HTTP para consultar:

- Lista de capitais;
- Lista de adjacências do grafo;
- Caminho de menor custo entre duas capitais utilizando o algoritmo de **Dijkstra**;
- Comparação entre a implementação própria do Dijkstra e a implementação da biblioteca `github.com/dominikbraun/graph`.

---

# Estrutura

```
backend/
├── api/
│   └── capitais.json          # Dataset das capitais
├── cmd/
│   └── main/
│       └── main.go            # Inicializa o servidor
├── docs/
│   ├── index.html             # Documentação da API
│   └── openapi.yaml
├── model/
│   └── city.go                # Modelo City
├── scripts/
│   └── makefile
├── server/
│   ├── handlers.go            # Endpoints HTTP
│   └── server.go              # Rotas e CORS
├── service/
│   ├── service.go             # Grafo + Dijkstra
│   └── service_test.go        # Testes automatizados
├── go.mod
├── go.sum
└── README.md
```

---

# Pré-requisitos

- Go 1.25+
- Biblioteca:

```
github.com/dominikbraun/graph
```

A biblioteca é utilizada **apenas para conferência dos resultados**.

O algoritmo principal de Dijkstra foi implementado manualmente utilizando:

- Lista de Adjacências
- Heap (Fila de Prioridade)
- `container/heap`

---

# Como executar

Na pasta `backend` execute:

```bash
go run ./cmd/main
```

ou

```bash
make run
```

Saída esperada:

```
grafo construído — 27 capitais
servidor iniciado em http://localhost:8081
```

---

# Como o grafo é construído

Cada capital representa um **vértice**.

Cada estrada entre capitais representa uma **aresta**.

O arquivo `capitais.json` possui o seguinte formato:

```json
{
    "Curitiba": {
        "toll": 30,
        "neighbors": {
            "São Paulo": 408,
            "Florianópolis": 300
        }
    }
}
```

Cada capital possui:

- pedágio (`toll`);
- lista de capitais vizinhas (`neighbors`);
- distância em quilômetros para cada vizinha.

Durante o carregamento do JSON, o método `Seed()` monta uma **Lista de Adjacências**, onde:

```
adj[A][B] = distância entre A e B
```

Como o grafo é **não direcionado**, toda aresta é criada nos dois sentidos.

Exemplo:

```
Curitiba -------- São Paulo
```

gera:

```
Curitiba -> São Paulo

São Paulo -> Curitiba
```

---

# Modelo de custo

O objetivo não é encontrar o menor caminho em quilômetros, mas sim o caminho de **menor custo financeiro**.

Cada trecho possui dois custos:

## Combustível

```
(distância / autonomia) × preço do litro
```

## Pedágio

Ao chegar em uma capital é cobrado o pedágio daquela cidade.

A capital de origem não paga pedágio.

Assim, o peso utilizado pelo algoritmo é:

```
peso(u,v)

=

combustível(u,v)

+

pedágio(v)
```

---

# Algoritmo

O sistema utiliza o algoritmo de **Dijkstra** com **Heap (Fila de Prioridade)**.

O algoritmo mantém:

- custo mínimo conhecido para cada capital;
- capital anterior (`prev`) para reconstrução do caminho;
- fila de prioridade contendo sempre a capital de menor custo acumulado.

Ao final são calculados:

- caminho encontrado;
- distância total;
- custo de combustível;
- custo dos pedágios;
- custo total.

---

# Endpoints

Todos retornam JSON.

---

## GET /capitais

Lista todas as capitais em ordem alfabética.

Exemplo:

```
GET /capitais
```

Resposta:

```json
{
  "count": 27,
  "capitais": [
    "Aracajú",
    "Belém",
    "Belo Horizonte"
  ]
}
```

---

## GET /show

Exibe a Lista de Adjacências do grafo.

```
GET /show
```

Resposta:

```json
{
  "count": 27,
  "adjacency": {
    "Curitiba": [
      {
        "name":"Florianópolis",
        "distance":300
      }
    ]
  }
}
```

---

## GET /caminho

Calcula o caminho de menor custo utilizando a implementação própria do algoritmo de Dijkstra.

Parâmetros:

```
origem
destino
combustivel
autonomia
```

Exemplo:

```
GET /caminho?origem=Curitiba&destino=Rio%20de%20Janeiro&combustivel=6&autonomia=12
```

Resposta:

```json
{
    "path":[
        "Curitiba",
        "São Paulo",
        "Rio de Janeiro"
    ],
    "distance":837,
    "fuelCost":418.5,
    "tollCost":70,
    "totalCost":488.5,
    "found":true
}
```

---

## GET /comparar

Executa dois algoritmos:

- implementação própria do Dijkstra;
- implementação da biblioteca `github.com/dominikbraun/graph`.

Permite verificar se ambos encontram o mesmo resultado.

Exemplo:

```
GET /comparar?origem=Curitiba&destino=Manaus&combustivel=6&autonomia=12
```

Resposta:

```json
{
    "meu": {...},
    "lib": {...},
    "custosBatem": true,
    "mesmoCaminho": true
}
```

---

# Testes

Os testes automatizados verificam:

- exemplo calculado manualmente;
- capitais sem rota;
- origem igual ao destino;
- comparação entre Dijkstra e Bellman-Ford;
- validação do caminho encontrado;
- conferência dos custos calculados.

Execute:

```bash
go test ./...
```

---

# Tecnologias utilizadas

- Go
- Lista de Adjacências
- Dijkstra
- Bellman-Ford (testes)
- Heap (`container/heap`)
- HTTP
- JSON
- OpenAPI
