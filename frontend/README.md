# Frontend · 8 Graus de Network

Interface em **Angular 21 + TypeScript** do trabalho **TD 01 — Teoria de Grafos** (UNESC, Ciência da Computação).

A aplicação consome o backend Go (`/actors`, `/bfs`, `/bfs8`, `/show`) e permite ao usuário:

- Selecionar um **Ator de Origem** e um **Ator de Destino** a partir de um `datalist` alimentado pela API `/actors`.
- Executar **BFS · Caminho mais curto** (`/bfs`), exibindo o caminho mínimo de vértices entre os atores e seu comprimento em arestas.
- Executar **BFS · Comprimento máximo 8** (`/bfs8`), exibindo **todos** os caminhos mínimos com no máximo 8 arestas.
- Visualizar o caminho como uma sequência de cartões, alternando entre vértices de **🎭 Ator** e **🎞️ Filme**.
- Tratamento explícito de erro quando não há relacionamento entre os atores (`Nenhum relacionamento encontrado`).

---

## Pré-requisitos

- **Node.js** ≥ 20
- **npm** ≥ 10
- O **backend** Go rodando em `http://localhost:8081` (ver `../backend`).

## Instalação

```bash
npm install
```

## Servidor de desenvolvimento

```bash
npm start
```

A aplicação fica disponível em `http://localhost:4200`.

> O `proxy.conf.json` redireciona automaticamente as rotas `/actors`, `/bfs`, `/bfs8`, `/show` e `/health` para `http://localhost:8081`, evitando problemas de CORS durante o desenvolvimento.

## Build de produção

```bash
npm run build
```

Os artefatos são gerados em `dist/frontend/`.

## Testes

```bash
npm test
```

Executa os testes unitários com [Vitest](https://vitest.dev/).

---

## Estrutura

```
frontend/
├─ proxy.conf.json          # Proxy para o backend em :8081
├─ src/
│  ├─ index.html
│  ├─ main.ts
│  ├─ styles.scss           # Tema base / variáveis CSS globais
│  └─ app/
│     ├─ app.config.ts      # Provê HttpClient para toda a aplicação
│     ├─ app.ts             # Componente raiz (formulário + resultados)
│     ├─ app.html           # Template (Angular control flow @if/@for)
│     ├─ app.scss           # Estilos do componente
│     └─ actor.service.ts   # Cliente HTTP do backend
```

## Endpoints consumidos

| Método | Rota                              | Descrição                                                |
|:------:|-----------------------------------|----------------------------------------------------------|
| GET    | `/actors`                         | Lista ordenada de todos os atores (alimenta o datalist). |
| GET    | `/bfs?from=A&to=B`                | BFS clássico — caminho mais curto entre dois atores.     |
| GET    | `/bfs8?from=A&to=B&max=8`         | Todos os caminhos mínimos com no máximo `max` arestas.   |
| GET    | `/show`                           | Mapa de adjacência completo (vértices e seus vizinhos).  |

---

## Tecnologias

- **Angular 21** standalone components com signals (`signal`, `computed`).
- **TypeScript 5.9**.
- Sintaxe nova de controle de fluxo no template (`@if`, `@for`, `@let`).
- **SCSS** com variáveis CSS (tema escuro).

---

UNESC · Ciência da Computação · Teoria de Grafos · TD 01
