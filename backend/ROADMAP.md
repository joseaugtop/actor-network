# Roadmap de estudo — 8 Graus de Network

Ordem sugerida para entender o código do zero até o algoritmo mais complexo.

---

## 1. Go básico (sintaxe)

Você não precisa dominar Go, só entender o suficiente para ler o código.

| Conceito | Onde aparece no código |
|---|---|
| `var` e `:=` | declaração de variáveis locais |
| `struct` | `model.Movie`, `service.Service` |
| `func` e múltiplos retornos | `ShortestPath` retorna `([]string, error)` |
| `for range` | loop sobre filmes e atores no `Seed` |
| `map[K]V` | `actors map[string]bool`, `parent map[string]string` |
| `slice` (`[]string`) | fila do BFS, lista de caminhos |
| `append` | adicionar elementos no slice |
| `slice[1:]` | remover o primeiro elemento (simular fila) |

**Referência rápida:** [tour.golang.org](https://go.dev/tour)

---

## 2. Grafos — conceitos fundamentais

Antes de olhar qualquer algoritmo, entenda a estrutura.

- **Vértice (nó):** cada elemento do grafo. Aqui: cada filme e cada ator.
- **Aresta:** ligação entre dois vértices. Aqui: liga um ator ao filme em que ele apareceu.
- **Grafo não direcionado:** a aresta não tem sentido — se A está ligado a B, B está ligado a A.
- **Grafo bipartido:** vértices divididos em dois grupos (filmes e atores) onde arestas só conectam grupos diferentes.
- **Lista de adjacências:** forma de representar o grafo em memória. Para cada vértice, guarda a lista dos seus vizinhos.

```
"Zendaya" → ["Dune: Part Two", "Spider-Man: Far From Home", ...]
"Dune: Part Two" → ["Zendaya", "Timothée Chalamet", "Rebecca Ferguson", ...]
```

**No código:** `service.AdjacencyMap()` — endpoint `/show`.

---

## 3. BFS — Busca em Largura

Usado em `ShortestPath` — endpoint `/bfs`.

### Ideia central
Visita o grafo **camada por camada**: primeiro os vizinhos diretos, depois os vizinhos dos vizinhos, e assim por diante. O primeiro caminho que chega no destino é o mais curto.

### Estruturas necessárias
- **Fila** (`queue`): guarda os vértices a visitar. Entra pelo fundo, sai pela frente (FIFO).
- **Mapa de pais** (`parent`): para cada vértice visitado, guarda de onde veio. Serve para reconstruir o caminho e também como "visitados".

### Passo a passo
```
Origem: "Zendaya"    Destino: "Tom Cruise"

fila:    [Zendaya]
parent:  {Zendaya: ""}

passo 1: tira Zendaya, adiciona seus filmes
fila:    [Dune 2, Spider-Man, ...]
parent:  {Zendaya:"", Dune 2:"Zendaya", Spider-Man:"Zendaya", ...}

passo 2: tira Dune 2, adiciona seus atores
fila:    [Spider-Man, ..., Timothée, Rebecca Ferguson, ...]
parent:  {..., Timothée:"Dune 2", Rebecca Ferguson:"Dune 2", ...}

... continua até achar "Tom Cruise" ...

achou! reconstrói: Tom Cruise → filme X → Rebecca Ferguson → Dune 2 → Zendaya
inverte: Zendaya → Dune 2 → Rebecca Ferguson → filme X → Tom Cruise
```

### Por que dá o caminho mínimo?
Porque a fila garante que você sempre processa os mais próximos primeiro. Quando o destino é encontrado, é impossível que exista um caminho mais curto — ele teria sido encontrado antes.

---

## 4. DFS — Busca em Profundidade

Base do `AllPathsUpTo` — endpoint `/bfs8`.

### Diferença do BFS
Enquanto BFS vai **por camadas**, DFS vai **fundo em um caminho até não ter mais onde ir**, depois volta (backtracking) e tenta outro.

### Backtracking
A chave do DFS para encontrar todos os caminhos simples:

```
visited["Brad Pitt"] = true
path = [..., "Brad Pitt"]

dfs(Brad Pitt, remaining-1)   // vai fundo

path = path[:len(path)-1]     // remove Brad Pitt do caminho atual
delete(visited, "Brad Pitt")  // marca como não visitado
                              // para outros caminhos poderem usar ele
```

Sem o `delete(visited, n)`, nenhum outro caminho poderia passar por "Brad Pitt".

---

## 5. Aprofundamento Iterativo (IDDFS)

É a técnica usada em `AllPathsUpTo` para enumerar **todos** os caminhos até tamanho 8.

### O problema de um DFS puro
Um DFS simples não garante ordenação por tamanho — poderia retornar um caminho de 7 arestas antes de um de 2.

### A solução
Executa DFS **várias vezes**, cada vez com um limite diferente:

```
rodada d=1: encontra todos os caminhos de exatamente 1 aresta
rodada d=2: encontra todos os caminhos de exatamente 2 arestas
rodada d=3: encontra todos os caminhos de exatamente 3 arestas
...
rodada d=8: encontra todos os caminhos de exatamente 8 arestas
```

Resultado: todos os caminhos menores são encontrados **antes** dos maiores. Se o limite de 10.000 caminhos for atingido, só os mais longos ficam de fora.

### No código

```go
for d := 1; d <= maxLen && !truncated; d++ {
    dfs(from, d)   // DFS que só aceita caminhos de exatamente d arestas
}
```

---

## 6. Ligando tudo ao código

| Conceito | Arquivo | Função |
|---|---|---|
| Seed do grafo | `service/service.go` | `Seed` |
| Lista de adjacências | `service/service.go` | `AdjacencyMap` |
| BFS caminho mínimo | `service/service.go` | `ShortestPath` |
| Reconstrução do caminho | `service/service.go` | `reconstruct` |
| DFS + aprofundamento iterativo | `service/service.go` | `AllPathsUpTo` |
| Endpoints HTTP | `server/handlers.go` | `handleBFS`, `handleBFS8` |
