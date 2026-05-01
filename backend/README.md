# 8 Graus de Network — Backend

Backend em Go para o trabalho **TD 01 — 8 Graus de Network** (Teoria de Grafos / UNESC).
Modela um grafo bipartido **Filme ↔ Ator** a partir de `api/latest_movies.json` e expõe uma
API HTTP para consultar os relacionamentos via BFS.

## Estrutura

```
backend/
├── api/latest_movies.json     # dataset de filmes/atores (seed)
├── cmd/main/main.go           # ponto de entrada — carrega o JSON e sobe o servidor
├── model/movie.go             # modelo Movie { id, title, cast }
├── service/service.go         # construção do grafo + algoritmos (BFS, BFS≤8)
├── server/server.go           # roteamento HTTP, handlers e CORS
├── scripts/makefile           # atalhos de build/run
├── go.mod / go.sum
└── README.md
```

## Pré-requisitos

- Go 1.25+
- A única dependência externa é [`github.com/dominikbraun/graph`](https://github.com/dominikbraun/graph),
  usada apenas como estrutura de armazenamento. **As buscas (BFS) são implementadas
  manualmente** em `service/service.go` para atender ao escopo do trabalho.

## Como executar

A partir da pasta `backend/`:

```bash
go run ./cmd/main
```

Ou via makefile (a partir de `backend/scripts/`):

```bash
make run
```

Saída esperada:

```
grafo construído — 1500 filmes, 8905 atores únicos
servidor iniciado em http://localhost:8081
```

## Como o grafo é construído

- **Vértices**: cada filme (`title`) e cada ator do `cast` viram um vértice (strings).
- **Arestas**: para cada par (filme, ator), são adicionadas arestas nos **dois sentidos**
  (filme → ator e ator → filme), conforme exigido pelo enunciado.
- Duplicatas no dataset (IDs repetidos, mesmo ator listado duas vezes em um filme) são
  toleradas: `AddVertex` ignora vértices repetidos e o erro de `AddEdge` para arestas já
  existentes é descartado.

## Endpoints

Todos os endpoints respondem JSON e suportam CORS (`Access-Control-Allow-Origin: *`).

### `GET /health`
Status do servidor.

```bash
curl http://localhost:8081/health
```

```json
{ "status": "ok", "mensagem": "servidor ativo" }
```

### `GET /actors`
Lista de todos os atores (ordenada alfabeticamente). Usada para popular `<select>`/`<datalist>`
no frontend.

```bash
curl http://localhost:8081/actors
```

```json
{
  "count": 8905,
  "actors": ["A. J. Cook", "Aaron Eckhart", "Aaron Paul", "..."]
}
```

### `GET /bfs?from=<ator>&to=<ator>`
Executa um **BFS clássico** e retorna o **caminho mínimo** entre dois atores.

```bash
curl "http://localhost:8081/bfs?from=Zendaya&to=Tom%20Cruise"
```

Sucesso:

```json
{
  "path": ["Zendaya", "The Greatest Showman", "Rebecca Ferguson", "Mission: Impossible - Dead Reckoning Part One", "Tom Cruise"],
  "length": 4,
  "found": true
}
```

Não encontrado:

```json
{ "path": [], "length": -1, "found": false, "message": "nenhum relacionamento encontrado" }
```

Ator inexistente (HTTP 404):

```json
{ "path": [], "length": -1, "found": false, "message": "ator não encontrado" }
```

### `GET /bfs8?from=<ator>&to=<ator>&max=<numero>`
Adaptação do BFS: retorna **todos os caminhos mínimos** cujo comprimento (em arestas)
seja ≤ `max` (padrão **8**). Útil quando há múltiplos relacionamentos curtos equivalentes.

```bash
curl "http://localhost:8081/bfs8?from=Zendaya&to=Samuel%20L.%20Jackson"
```

```json
{
  "paths": [
    {
      "path": ["Zendaya", "Spider-Man: Far From Home", "Samuel L. Jackson"],
      "length": 2
    }
  ],
  "count": 1,
  "length": 2,
  "found": true
}
```

Quando não há relação dentro do limite:

```json
{ "paths": [], "count": 0, "found": false, "message": "nenhum relacionamento encontrado dentro do comprimento máximo" }
```

## Algoritmos

- **BFS (`/bfs`)**: BFS padrão sobre lista de adjacência, com mapa de pais para reconstruir
  o caminho mínimo.
- **BFS adaptado (`/bfs8`)**: BFS que registra **todos** os predecessores no mesmo nível
  mínimo (DAG de predecessores). Após encontrar o destino, faz uma travessia recursiva
  do destino até a origem para enumerar **todos** os caminhos de comprimento mínimo
  ≤ `max`. A poda no nível do alvo garante terminação rápida mesmo em grafos densos.

> Observação: como o grafo é bipartido, todo caminho entre atores tem **comprimento par**
> (alterna ator → filme → ator → ...).

## Notas sobre o dataset

O `latest_movies.json` contém **1500 filmes** e **8905 atores únicos**, com algumas
duplicatas conhecidas:

- ~25 IDs de filmes repetidos
- ~50 títulos de filmes repetidos (remakes ou entradas redundantes)
- 4 filmes com o mesmo ator listado duas vezes no `cast`

Essas duplicatas não corrompem o grafo: o builder as absorve silenciosamente.
