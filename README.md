# Caminho Mais Barato entre Capitais

Trabalho de Desenvolvimento 02 (TD 02) da disciplina de **Teoria de Grafos** do curso de **Ciência da Computação** da [UNESC](https://www.unesc.net/), sob orientação do **Prof. André Faria Ruaro**.

---

# Sobre o projeto

Aplicação web que encontra o **caminho de menor custo entre duas capitais brasileiras** utilizando o algoritmo de **Dijkstra**.

As capitais são carregadas a partir do arquivo `capitais.json`, formando um **grafo não direcionado representado por Lista de Adjacências**.

Cada capital representa um vértice e cada estrada representa uma aresta com sua respectiva distância em quilômetros. Além da distância, cada capital possui um valor de pedágio utilizado no cálculo do custo total da viagem.

O custo de cada trecho é calculado considerando:

- consumo de combustível;
- preço do combustível;
- pedágios das capitais visitadas.

---

# Funcionalidades

- Carregamento automático do grafo a partir do arquivo `capitais.json`;
- Visualização da Lista de Adjacências;
- Listagem de todas as capitais cadastradas;
- Cálculo do caminho de menor custo utilizando o algoritmo de Dijkstra;
- Comparação entre a implementação própria do algoritmo e a biblioteca `dominikbraun/graph`;
- Interface web para consulta das rotas;
- API REST documentada.

---

# Algoritmos utilizados

- Lista de Adjacências
- Dijkstra (implementação própria utilizando Heap/Fila de Prioridade)
- Dijkstra da biblioteca `github.com/dominikbraun/graph` (utilizado apenas para conferência dos resultados)
- Bellman-Ford (utilizado nos testes automatizados)

---

# Tecnologias

## Backend

- Go
- HTTP
- JSON
- OpenAPI
- container/heap
- github.com/dominikbraun/graph

## Frontend

- Angular
- TypeScript
- HTML
- SCSS

---

# Estrutura do projeto

```
projeto/
├── backend/
│   ├── api/
│   │   └── capitais.json
│   ├── cmd/
│   │   └── main/
│   ├── docs/
│   ├── model/
│   ├── server/
│   ├── service/
│   ├── scripts/
│   ├── go.mod
│   └── README.md
│
└── frontend/
    ├── src/
    ├── public/
    ├── angular.json
    ├── package.json
    └── README.md
```

---

# Modelo do Grafo

- **Vértices:** Capitais brasileiras
- **Arestas:** Rodovias entre capitais
- **Peso:** Custo da viagem

O peso utilizado pelo algoritmo é calculado pela fórmula:

```
peso = combustível + pedágio
```

onde:

```
combustível = (distância ÷ autonomia) × preço do litro
```

O pedágio é cobrado ao chegar em cada capital (exceto a origem).

---

# Funcionalidades da API

O backend disponibiliza os seguintes endpoints:

| Método | Endpoint | Descrição |
|---------|----------|-----------|
| GET | `/capitais` | Lista todas as capitais |
| GET | `/show` | Exibe a Lista de Adjacências |
| GET | `/caminho` | Calcula o caminho de menor custo |
| GET | `/comparar` | Compara o Dijkstra implementado com a biblioteca |

---

# Testes

O projeto possui testes automatizados que verificam:

- exemplo calculado manualmente;
- capitais sem conexão;
- origem igual ao destino;
- comparação entre Dijkstra e Bellman-Ford;
- validação dos caminhos encontrados.

---

# Disciplina

| Campo | Informação |
|--------|------------|
| Curso | Ciência da Computação |
| Disciplina | Teoria de Grafos |
| Trabalho | TD 02 — Caminho Mais Barato entre Capitais |
| Professor | André Faria Ruaro |
| Instituição | UNESC |
