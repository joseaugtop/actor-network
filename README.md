# 8 Graus de Network

Trabalho de Desenvolvimento 01 (TD 01) da disciplina de **Teoria de Grafos** do curso de **Ciência da Computação** da [UNESC](https://www.unesc.net/), sob orientação do **Prof. André Faria Ruaro**.

---

## Sobre o projeto

Aplicação que encontra o relacionamento mais próximo entre dois atores de cinema.

Os dados são carregados a partir de um arquivo `latest_movies.json` e estruturados em um **Grafo não direcionado de Lista de Adjacências**, onde cada ator e cada filme são vértices, e a participação de um ator em um filme representa uma aresta.

---

## Funcionalidades

- **Seed** — carrega os dados do JSON e monta o grafo
- **Show** — exibe todos os vértices e seus adjacentes
- **BFS** — encontra o caminho mínimo entre dois atores
- **BFS com limite de 8** — encontra todos os caminhos com comprimento máximo de 8 arestas entre os atores selecionados
- Interface com seleção de ator de origem e destino
- Tratamento de caso sem relacionamento encontrado

---

## Tecnologias

- **Backend** — Go + [dominikbraun/graph](https://github.com/dominikbraun/graph)
- **Frontend** — React

---

## Estrutura do projeto

```
actor-network/
├── backend/
│   ├── api/
│   │   └── latest_movies.json
│   ├── cmd/main/
│   │   └── main.go
│   ├── pkg/
│   │   └── movie.go
│   └── go.mod
└── frontend/              # em desenvolvimento
```

---

## Disciplina

| Campo | Info |
|---|---|
| Curso | Ciência da Computação |
| Disciplina | Teoria de Grafos |
| Trabalho | TD 01 — 8 Graus de Network |
| Professor | André Faria Ruaro |
| Instituição | UNESC |
