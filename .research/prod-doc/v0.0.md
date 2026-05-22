# Fintrago — Product Document
**FIN**ancial **TRA**cker with RA**G** and **Go**

> **Authors:** Shrey & Reece
> **Status:** Pre-development
> **Primary goal:** A hands-on portfolio project to prepare for IBM. Each person goes deep on their own stack. The project is intentionally minimal — real enough to be meaningful, scoped tightly enough to actually finish.

---

## 1. What Is Fintrago?

Fintrago is a personal finance tracker with an AI layer. Users manually log their income and expenses, upload financial documents, and query their data in plain English using a RAG-powered AI agent.

The project is split across two services — a Go REST API and a Python AI service — so that Reece and Shrey can each go deep on their respective stacks while building something that works end to end.

**This is a learning project first.** Every technology choice is made because one of us needs to learn it, not because it's the fanciest option.

---

## 2. Core Features (MVP)

1. **Manual transaction management** — add, view, and delete income/expense entries via REST API
2. **Balance summary** — total balance and basic spending breakdown computed from stored transactions
3. **Document upload** — accept a PDF or CSV, chunk and embed it, store in a vector database
4. **Natural language Q&A** — ask questions in plain English; a LangGraph agent retrieves relevant context and answers using an LLM

That's it for MVP. No frontend, no bank sync, no multi-user. All interaction happens via Postman or curl.

---

## 3. Who Does What

| Area | Owner | Stack |
|---|---|---|
| REST API (all routes) | Reece | Go, Gin |
| Database schema & migrations | Reece | PostgreSQL, pgx |
| Transaction CRUD | Reece | Go, Gin, Postgres |
| Balance & summary endpoints | Reece | Go, Gin, Postgres |
| Auth middleware (JWT) | Reece | Go, Gin |
| Request logging middleware | Reece | Go, Gin |
| Dockerfile (go-api) | Reece | Docker |
| `/upload` — receives file, proxies to ai-service | Reece | Go, Gin |
| `/ask` — receives question, proxies to ai-service | Reece | Go, Gin |
| FastAPI service skeleton | Shrey | Python, FastAPI |
| Document loading & chunking | Shrey | Python, LangChain |
| Embedding & vector storage | Shrey | LangChain, ChromaDB |
| LangGraph agent + tools | Shrey | Python, LangGraph |
| RAG retrieval pipeline | Shrey | Python, LangChain |
| Dockerfile (ai-service) | Shrey | Docker |
| docker-compose.yml (full stack) | Shrey | Docker Compose |
| GitHub Actions CI pipeline | Shrey | GitHub Actions |
| OpenShift deployment | Shrey | RedHat OpenShift |

---

## 4. File Structure

```
fintrago/
│
├── go-api/                               # Reece
│   ├── main.go
│   ├── config/
│   │   └── config.go
│   ├── routes/
│   │   └── routes.go
│   ├── handlers/
│   │   ├── transactions.go
│   │   ├── balance.go
│   │   └── proxy.go
│   ├── models/
│   │   └── transaction.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── logger.go
│   ├── db/
│   │   ├── db.go
│   │   └── migrations/
│   │       └── 001_init.sql
│   ├── Dockerfile
│   └── go.mod / go.sum
│
├── ai-service/                           # Shrey
│   ├── main.py
│   ├── config.py
│   ├── routers/
│   │   ├── upload.py
│   │   └── ask.py
│   ├── rag/
│   │   ├── loader.py
│   │   ├── embedder.py
│   │   ├── retriever.py
│   │   └── prompt_builder.py
│   ├── agents/
│   │   ├── graph.py
│   │   ├── nodes.py
│   │   └── tools.py
│   ├── db/
│   │   └── vector_store.py
│   ├── tests/
│   │   ├── test_rag.py
│   │   └── test_agent.py
│   ├── requirements.txt
│   └── Dockerfile
│
├── docker-compose.yml
├── docker-compose.override.yml
└── .github/
    └── workflows/
        ├── ci.yml
        └── deploy.yml
```

---

## 5. Shared JSON Contract

Lock this in before either side starts building. This is the only interface between the two services.

```json
// POST /ask — request
{ "user_id": "uuid", "question": "How much did I spend on food in April?" }

// POST /ask — response
{
  "answer": "You spent $340 on food in April across 12 transactions.",
  "sources": [
    { "type": "transaction", "id": "uuid", "description": "Whole Foods", "amount": 54.20 },
    { "type": "document", "filename": "april_statement.pdf", "chunk": "..." }
  ]
}

// POST /upload — request (multipart/form-data)
{ "user_id": "uuid", "file": "<binary PDF or CSV>" }

// POST /upload — response
{ "message": "Uploaded and indexed successfully.", "filename": "april_statement.pdf", "chunks_stored": 42 }
```

---

## 6. Reece — Step-by-Step Build Guide

### Phase 1: Project scaffold

**Step 1 — Initialise the Go module**
```bash
mkdir fintrago && cd fintrago
mkdir go-api && cd go-api
go mod init github.com/yourname/fintrago
go get github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv
```

**Step 2 — Create the folder structure**
```bash
mkdir -p config routes handlers models middleware db/migrations
touch main.go config/config.go routes/routes.go \
      handlers/transactions.go handlers/balance.go handlers/proxy.go \
      models/transaction.go middleware/auth.go middleware/logger.go \
      db/db.go db/migrations/001_init.sql
```

**Step 3 — Write `main.go`**

This is the entry point. It loads config, connects to the DB, registers routes, and starts the server.

```go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/yourname/fintrago/config"
    "github.com/yourname/fintrago/db"
    "github.com/yourname/fintrago/routes"
)

func main() {
    cfg := config.Load()
    db.Connect(cfg.DatabaseURL)

    r := gin.Default()
    routes.Register(r)

    log.Printf("Server running on port %s", cfg.Port)
    r.Run(":" + cfg.Port)
}
```

**Step 4 — Write `config/config.go`**

Reads environment variables from `.env`. Never hardcode secrets.

```go
package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port          string
    DatabaseURL   string
    JWTSecret     string
    AIServiceURL  string
}

func Load() *Config {
    godotenv.Load()
    return &Config{
        Port:         getEnv("PORT", "8080"),
        DatabaseURL:  getEnv("DATABASE_URL", ""),
        JWTSecret:    getEnv("JWT_SECRET", ""),
        AIServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:8000"),
    }
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}
```

---

### Phase 2: Data model and in-memory CRUD

**Step 5 — Define the Transaction model**

```go
// models/transaction.go
package models

import "time"

type TransactionType string

const (
    Income  TransactionType = "income"
    Expense TransactionType = "expense"
)

type Transaction struct {
    ID          string          `json:"id"`
    UserID      string          `json:"user_id"`
    Type        TransactionType `json:"type"`
    Amount      float64         `json:"amount"`
    Category    string          `json:"category"`
    Description string          `json:"description"`
    Date        time.Time       `json:"date"`
    CreatedAt   time.Time       `json:"created_at"`
}
```

**Step 6 — Write in-memory handlers first (no DB yet)**

Start with a slice as the data store. This lets you test routing and JSON handling before touching Postgres.

```go
// handlers/transactions.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/yourname/fintrago/models"
)

var store []models.Transaction  // temporary in-memory store

func CreateTransaction(c *gin.Context) {
    var t models.Transaction
    if err := c.ShouldBindJSON(&t); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    store = append(store, t)
    c.JSON(http.StatusCreated, t)
}

func ListTransactions(c *gin.Context) {
    c.JSON(http.StatusOK, store)
}

func DeleteTransaction(c *gin.Context) {
    id := c.Param("id")
    for i, t := range store {
        if t.ID == id {
            store = append(store[:i], store[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "deleted"})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}
```

**Step 7 — Register routes**

```go
// routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/yourname/fintrago/handlers"
)

func Register(r *gin.Engine) {
    api := r.Group("/api")
    {
        api.POST("/transactions", handlers.CreateTransaction)
        api.GET("/transactions", handlers.ListTransactions)
        api.DELETE("/transactions/:id", handlers.DeleteTransaction)
        api.GET("/balance", handlers.GetBalance)
        api.POST("/upload", handlers.ProxyUpload)
        api.POST("/ask", handlers.ProxyAsk)
    }
}
```

Test all routes with Postman before moving on. Make sure JSON in and out works correctly.

---

### Phase 3: PostgreSQL

**Step 8 — Write the migration SQL**

```sql
-- db/migrations/001_init.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS transactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     TEXT NOT NULL,
    type        TEXT CHECK (type IN ('income', 'expense')) NOT NULL,
    amount      NUMERIC(12, 2) NOT NULL,
    category    TEXT,
    description TEXT,
    date        DATE NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
```

**Step 9 — Write `db/db.go`**

```go
package db

import (
    "context"
    "log"
    "github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(databaseURL string) {
    var err error
    Pool, err = pgxpool.New(context.Background(), databaseURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    log.Println("Connected to PostgreSQL")
}
```

**Step 10 — Replace in-memory store with Postgres**

Rewrite `handlers/transactions.go` to use `db.Pool` instead of the slice. Use `pgx` to run `INSERT`, `SELECT`, and `DELETE` queries. Scan rows into `models.Transaction` structs.

Key pgx patterns to know:
- `db.Pool.QueryRow(ctx, "SELECT ...", args...)` — one row
- `db.Pool.Query(ctx, "SELECT ...", args...)` — multiple rows
- `db.Pool.Exec(ctx, "INSERT ...", args...)` — write
- `rows.Scan(&t.ID, &t.Amount, ...)` — map columns to struct fields

---

### Phase 4: Remaining endpoints

**Step 11 — Write `handlers/balance.go`**

Query Postgres to sum income and expenses separately, return the net balance plus a breakdown by category.

**Step 12 — Write `handlers/proxy.go`**

These routes receive the request from the client and forward it to Shrey's Python service. Use Go's `net/http` to make an outbound HTTP request to `AI_SERVICE_URL`.

```go
// Example: forward POST /ask to ai-service
func ProxyAsk(c *gin.Context) {
    body, _ := io.ReadAll(c.Request.Body)
    resp, err := http.Post(aiServiceURL+"/ask", "application/json", bytes.NewReader(body))
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "ai-service unreachable"})
        return
    }
    defer resp.Body.Close()
    result, _ := io.ReadAll(resp.Body)
    c.Data(resp.StatusCode, "application/json", result)
}
```

---

### Phase 5: Middleware and auth

**Step 13 — Write `middleware/logger.go`**

Gin has a built-in logger (`gin.Logger()`). Write a custom one that logs the method, path, status code, and latency. Add it to the router with `r.Use(middleware.Logger())`.

**Step 14 — Write `middleware/auth.go`**

For MVP, keep auth simple: require a Bearer JWT token on all `/api` routes. Validate the token using `golang-jwt`. If invalid, return 401. If valid, extract `user_id` from claims and set it on the Gin context so handlers can use it.

---

### Phase 6: Docker

**Step 15 — Write the Dockerfile**

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o fintrago-api ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fintrago-api .
EXPOSE 8080
CMD ["./fintrago-api"]
```

Use a multi-stage build: the first stage compiles the binary, the second stage is a tiny Alpine image that just runs it. The final image is under 20MB.

**Step 16 — Test the Docker build**
```bash
docker build -t fintrago-api .
docker run -p 8080:8080 --env-file .env fintrago-api
```

---

## 7. Shrey — Step-by-Step Build Guide

### Phase 1: Project scaffold

**Step 1 — Initialise the Python project**
```bash
cd fintrago
mkdir ai-service && cd ai-service
python -m venv venv && source venv/bin/activate
pip install fastapi uvicorn langchain langchain-openai \
            langchain-community chromadb pypdf python-multipart \
            python-dotenv pytest httpx
pip freeze > requirements.txt
```

**Step 2 — Create the folder structure**
```bash
mkdir -p routers rag agents db tests
touch main.py config.py \
      routers/upload.py routers/ask.py \
      rag/loader.py rag/embedder.py rag/retriever.py rag/prompt_builder.py \
      agents/graph.py agents/nodes.py agents/tools.py \
      db/vector_store.py \
      tests/test_rag.py tests/test_agent.py
```

**Step 3 — Write `config.py`**

```python
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    openai_api_key: str
    chroma_host: str = "chromadb"
    chroma_port: int = 8001
    port: int = 8000
    langchain_tracing_v2: bool = False
    langchain_api_key: str = ""

    class Config:
        env_file = ".env"

settings = Settings()
```

**Step 4 — Write `main.py`**

```python
from fastapi import FastAPI
from routers import upload, ask

app = FastAPI(title="Fintrago AI Service")

app.include_router(upload.router, prefix="/upload", tags=["upload"])
app.include_router(ask.router, prefix="/ask", tags=["ask"])

@app.get("/health")
def health():
    return {"status": "ok"}
```

Run it:
```bash
uvicorn main:app --reload --port 8000
```

Visit `http://localhost:8000/docs` — FastAPI auto-generates interactive API docs. Use this instead of Postman for testing your own service.

---

### Phase 2: Document upload and ingestion pipeline

This is the ingest side of RAG: load → chunk → embed → store.

**Step 5 — Write `rag/loader.py`**

Handles loading raw documents and splitting them into chunks.

```python
from langchain_community.document_loaders import PyPDFLoader, CSVLoader
from langchain.text_splitter import RecursiveCharacterTextSplitter

def load_and_chunk(file_path: str, file_type: str) -> list:
    if file_type == "pdf":
        loader = PyPDFLoader(file_path)
    elif file_type == "csv":
        loader = CSVLoader(file_path)
    else:
        raise ValueError(f"Unsupported file type: {file_type}")

    documents = loader.load()

    splitter = RecursiveCharacterTextSplitter(
        chunk_size=500,
        chunk_overlap=50,
        separators=["\n\n", "\n", ".", " "]
    )
    return splitter.split_documents(documents)
```

Chunking strategy explained:
- `chunk_size=500` — each chunk is ~500 characters. Small enough to be specific, large enough to have context.
- `chunk_overlap=50` — chunks overlap by 50 characters so context isn't lost at boundaries.
- `RecursiveCharacterTextSplitter` tries to split on paragraph breaks first, then sentences, then words. This preserves meaning better than splitting at a fixed character count.

**Step 6 — Write `db/vector_store.py`**

```python
import chromadb
from config import settings

client = chromadb.HttpClient(host=settings.chroma_host, port=settings.chroma_port)

def get_collection(user_id: str):
    # Each user gets their own ChromaDB collection
    return client.get_or_create_collection(name=f"user_{user_id}")

def store_chunks(user_id: str, chunks: list, embeddings: list, metadatas: list):
    collection = get_collection(user_id)
    ids = [f"{user_id}_{i}" for i in range(len(chunks))]
    collection.add(documents=chunks, embeddings=embeddings, metadatas=metadatas, ids=ids)

def query(user_id: str, query_embedding: list, n_results: int = 5):
    collection = get_collection(user_id)
    return collection.query(query_embeddings=[query_embedding], n_results=n_results)
```

**Step 7 — Write `rag/embedder.py`**

```python
from langchain_openai import OpenAIEmbeddings

embedder = OpenAIEmbeddings(model="text-embedding-3-small")

def embed_texts(texts: list[str]) -> list[list[float]]:
    return embedder.embed_documents(texts)

def embed_query(query: str) -> list[float]:
    return embedder.embed_query(query)
```

`text-embedding-3-small` produces 1536-dimensional vectors. It's fast and cheap — good for a learning project.

**Step 8 — Write `routers/upload.py`**

```python
import shutil, tempfile, os
from fastapi import APIRouter, UploadFile, Form
from rag.loader import load_and_chunk
from rag.embedder import embed_texts
from db.vector_store import store_chunks

router = APIRouter()

@router.post("")
async def upload_document(user_id: str = Form(...), file: UploadFile = ...):
    suffix = ".pdf" if file.content_type == "application/pdf" else ".csv"

    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        shutil.copyfileobj(file.file, tmp)
        tmp_path = tmp.name

    file_type = "pdf" if suffix == ".pdf" else "csv"
    chunks = load_and_chunk(tmp_path, file_type)
    os.unlink(tmp_path)

    texts = [c.page_content for c in chunks]
    embeddings = embed_texts(texts)
    metadatas = [{"filename": file.filename, "chunk_index": i} for i in range(len(chunks))]

    store_chunks(user_id, texts, embeddings, metadatas)

    return {"message": "Uploaded and indexed successfully.", "filename": file.filename, "chunks_stored": len(chunks)}
```

Test this end to end before moving on: upload a PDF, query ChromaDB directly, verify chunks are stored.

---

### Phase 3: Retrieval

**Step 9 — Write `rag/retriever.py`**

```python
from rag.embedder import embed_query
from db.vector_store import query

def retrieve(user_id: str, question: str, n_results: int = 5) -> list[dict]:
    q_embedding = embed_query(question)
    results = query(user_id, q_embedding, n_results=n_results)

    chunks = []
    for i, doc in enumerate(results["documents"][0]):
        chunks.append({
            "content": doc,
            "metadata": results["metadatas"][0][i],
            "distance": results["distances"][0][i]
        })
    return chunks
```

**Step 10 — Write `rag/prompt_builder.py`**

This assembles the retrieved chunks into a prompt for the LLM.

```python
def build_prompt(question: str, context_chunks: list[dict]) -> str:
    context = "\n\n".join([
        f"[Source: {c['metadata'].get('filename', 'unknown')}, chunk {c['metadata'].get('chunk_index', '?')}]\n{c['content']}"
        for c in context_chunks
    ])

    return f"""You are a personal finance assistant. Answer the user's question using only the context below.
If the context doesn't contain enough information to answer, say so clearly.

Context:
{context}

Question: {question}

Answer:"""
```

---

### Phase 4: LangGraph agent

This is the most conceptually complex part. Take time to understand LangGraph before coding.

**Concepts to understand first:**
- A LangGraph agent is a directed graph where each node is a function
- The graph has a shared `State` dict that each node reads from and writes to
- Edges between nodes can be conditional — the agent decides where to go next based on state
- Tools are functions the agent can call to interact with the outside world

**Step 11 — Define the agent's tools in `agents/tools.py`**

```python
from langchain_core.tools import tool
from rag.retriever import retrieve

@tool
def search_documents(user_id: str, query: str) -> str:
    """Search the user's uploaded financial documents for relevant information."""
    chunks = retrieve(user_id, query)
    return "\n\n".join([c["content"] for c in chunks])

# Future tool: query_transactions (calls Reece's /transactions endpoint)
# For MVP, this can be a stub that returns placeholder data
@tool
def query_transactions(user_id: str, query: str) -> str:
    """Query the user's transaction history."""
    return "Transaction querying not yet implemented. Use uploaded documents."
```

**Step 12 — Define the graph in `agents/graph.py`**

```python
from langgraph.graph import StateGraph, END
from langchain_openai import ChatOpenAI
from agents.nodes import call_model, call_tools, should_continue
from agents.tools import search_documents, query_transactions
import operator
from typing import Annotated, TypedDict

class AgentState(TypedDict):
    user_id: str
    question: str
    messages: Annotated[list, operator.add]

tools = [search_documents, query_transactions]
llm = ChatOpenAI(model="gpt-4o", temperature=0).bind_tools(tools)

def build_graph():
    graph = StateGraph(AgentState)
    graph.add_node("agent", call_model)
    graph.add_node("tools", call_tools)
    graph.set_entry_point("agent")
    graph.add_conditional_edges("agent", should_continue, {"tools": "tools", "end": END})
    graph.add_edge("tools", "agent")
    return graph.compile()
```

**Step 13 — Write `agents/nodes.py`**

```python
from langchain_core.messages import HumanMessage, AIMessage, ToolMessage
from langchain_openai import ChatOpenAI
from agents.tools import search_documents, query_transactions

tools_map = {
    "search_documents": search_documents,
    "query_transactions": query_transactions
}

llm = ChatOpenAI(model="gpt-4o", temperature=0).bind_tools(list(tools_map.values()))

def call_model(state):
    response = llm.invoke(state["messages"])
    return {"messages": [response]}

def call_tools(state):
    last_message = state["messages"][-1]
    tool_results = []
    for tool_call in last_message.tool_calls:
        tool = tools_map[tool_call["name"]]
        result = tool.invoke({**tool_call["args"], "user_id": state["user_id"]})
        tool_results.append(ToolMessage(content=result, tool_call_id=tool_call["id"]))
    return {"messages": tool_results}

def should_continue(state):
    last = state["messages"][-1]
    if hasattr(last, "tool_calls") and last.tool_calls:
        return "tools"
    return "end"
```

**Step 14 — Write `routers/ask.py`**

```python
from fastapi import APIRouter
from pydantic import BaseModel
from langchain_core.messages import HumanMessage
from agents.graph import build_graph

router = APIRouter()
graph = build_graph()

class AskRequest(BaseModel):
    user_id: str
    question: str

@router.post("")
async def ask(request: AskRequest):
    result = graph.invoke({
        "user_id": request.user_id,
        "question": request.question,
        "messages": [HumanMessage(content=request.question)]
    })
    final_message = result["messages"][-1]
    return {"answer": final_message.content, "sources": []}
```

---

### Phase 5: Docker and Compose

**Step 15 — Write the ai-service Dockerfile**

```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

**Step 16 — Write `docker-compose.yml` (at repo root)**

This is yours to own. It wires all four containers together: go-api, ai-service, postgres, chromadb.

```yaml
version: "3.9"

services:
  go-api:
    build: ./go-api
    ports:
      - "8080:8080"
    env_file: ./go-api/.env
    depends_on:
      - postgres

  ai-service:
    build: ./ai-service
    ports:
      - "8000:8000"
    env_file: ./ai-service/.env
    depends_on:
      - chromadb

  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: fintrago
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./go-api/db/migrations:/docker-entrypoint-initdb.d

  chromadb:
    image: chromadb/chroma:latest
    ports:
      - "8001:8001"
    volumes:
      - chromadata:/chroma/chroma

volumes:
  pgdata:
  chromadata:
```

Test the full stack:
```bash
docker compose up --build
```

All four containers should start. Hit `http://localhost:8080/api/transactions` — it should return an empty array from Postgres. Hit `http://localhost:8000/health` — should return `{"status": "ok"}`.

---

### Phase 6: CI/CD and OpenShift

**Step 17 — Write `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test-go:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: cd go-api && go test ./...

  test-python:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
      - run: cd ai-service && pip install -r requirements.txt && pytest

  build-images:
    runs-on: ubuntu-latest
    needs: [test-go, test-python]
    steps:
      - uses: actions/checkout@v4
      - run: docker build ./go-api -t fintrago-api
      - run: docker build ./ai-service -t fintrago-ai
```

**Step 18 — OpenShift deployment**

OpenShift is Kubernetes with extra conventions. For each service you need three YAML manifests:

- `Deployment` — what container to run, how many replicas, env vars
- `Service` — internal networking (so go-api can reach ai-service by name)
- `Route` — external URL (OpenShift-specific, replaces Kubernetes Ingress)

Use the RedHat Developer Sandbox (free) to deploy. Key CLI commands:
```bash
oc login --token=<your-token> --server=<your-cluster-url>
oc new-project fintrago
oc apply -f openshift/go-api-deployment.yaml
oc apply -f openshift/ai-service-deployment.yaml
oc get pods     # verify they're running
oc get routes   # get the external URLs
```

Secrets (API keys, DB passwords) go in OpenShift Secrets, not in the manifests directly:
```bash
oc create secret generic fintrago-secrets \
  --from-literal=openai-api-key=$OPENAI_API_KEY \
  --from-literal=jwt-secret=$JWT_SECRET
```

---

## 8. MVP Checklist

The MVP is complete when this full flow works end to end from a cold `docker compose up`:

- [ ] `POST /api/transactions` — creates a transaction, stored in Postgres
- [ ] `GET /api/transactions` — returns the list from Postgres
- [ ] `DELETE /api/transactions/:id` — removes from Postgres
- [ ] `GET /api/balance` — returns correct net balance
- [ ] `POST /api/upload` — accepts a PDF, chunks and embeds it, returns `chunks_stored`
- [ ] `POST /api/ask` — returns a coherent, grounded answer to a financial question
- [ ] `docker compose up --build` — all four services start cleanly with one command
- [ ] GitHub Actions CI — runs and passes on every push to main

---

## 9. Future Context — For Claude Agent Handoff

> **This section is a self-contained briefing for a future AI agent (or Claude conversation) to pick up where this document leaves off. Paste this section as context when starting that conversation.**

---

### What Fintrago is

Fintrago (FINancial TRAcker with RAG and Go) is a personal finance tracker with an AI Q&A layer. It is a portfolio/learning project built by two developers preparing for work at IBM.

**Reece** owns the Go/Gin REST API and PostgreSQL layer.
**Shrey** owns the Python/FastAPI AI service, LangChain/LangGraph RAG pipeline, Docker, GitHub Actions CI, and OpenShift deployment.

All client interaction is via REST API (Postman/curl) — no frontend exists yet.

### What was built in the MVP

- Go/Gin REST API: transaction CRUD, balance endpoint, JWT auth, proxy routes to ai-service
- Python/FastAPI AI service: document upload + chunking + embedding (OpenAI), ChromaDB vector store, LangGraph agent with `search_documents` and `query_transactions` tools, `/ask` endpoint
- PostgreSQL for structured transaction data
- ChromaDB for vector embeddings
- docker-compose wiring all four services: go-api, ai-service, postgres, chromadb
- GitHub Actions CI: lint + test + build on every push

### What comes next (post-MVP features)

**1. Plaid integration (Reece)**
The biggest planned feature. Connects real bank accounts so transactions are synced automatically rather than entered manually.
- Use Plaid Sandbox first (no real bank needed)
- New endpoints: `POST /plaid/link-token`, `POST /plaid/exchange-token`, `POST /plaid/sync`, `POST /plaid/webhook`
- Synced transactions go into the same Postgres `transactions` table
- After sync, new transactions should also be embedded and stored in ChromaDB (Shrey needs to expose an internal `/embed-transactions` endpoint, or Reece calls `/upload` with a CSV of new transactions)
- Go Plaid SDK: `github.com/plaid/plaid-go`

**2. `query_transactions` tool (Shrey)**
Currently a stub. Should make an HTTP call to Reece's `GET /api/transactions` endpoint (with filters: date range, category) and return results as structured text for the LLM to reason over. This is what enables questions like "how much did I spend on food last month?" to use live Postgres data, not just uploaded docs.

**3. RAG optimization (Shrey)**
- Experiment with chunk sizes (try 200, 500, 1000 — evaluate retrieval quality)
- Add MMR (Maximal Marginal Relevance) retrieval to reduce redundant chunks
- Add a cross-encoder reranker (e.g. `cross-encoder/ms-marco-MiniLM-L-6-v2`) to rerank top-k results before passing to LLM
- Switch from ChromaDB to pgvector for production (simplifies OpenShift deployment — one fewer service)
- Add LangSmith tracing for debugging agent runs

**4. OpenShift deployment (Shrey)**
- Write Deployment + Service + Route manifests for go-api and ai-service
- Store secrets in OpenShift Secrets (not ConfigMaps)
- Wire deploy.yml GitHub Actions workflow to push images to Quay.io and apply manifests via `oc`
- Consider using OpenShift's BuildConfig instead of GitHub Actions for image builds

**5. Spending analytics endpoints (Reece)**
- `GET /api/analytics/by-category?month=2025-04` — spending grouped by category
- `GET /api/analytics/trends?months=3` — month-over-month comparison
- These feed the agent's `query_transactions` tool with richer structured data

**6. Frontend (both)**
- Simple React or Next.js dashboard
- Pages: transaction list, balance summary, chat interface for `/ask`
- The chat interface is the most valuable — it makes the RAG pipeline visible to a real user

**7. Multi-user support (both)**
- Currently `user_id` is passed in every request body — no real auth
- Proper auth: JWT issued on login, `user_id` extracted from token in middleware (not from request body)
- Vector store already namespaced by `user_id` in ChromaDB collections — no change needed there
- Postgres `transactions` table already has `user_id` column — add index on it

### Tech stack reference

| Service | Language | Framework | Key libraries |
|---|---|---|---|
| go-api | Go 1.22 | Gin | pgx v5, golang-jwt, godotenv |
| ai-service | Python 3.11 | FastAPI | LangChain, LangGraph, OpenAI, ChromaDB, pypdf |
| Database | — | PostgreSQL 15 | pgvector extension (for post-MVP) |
| Vector store | — | ChromaDB (dev) | pgvector (prod, post-MVP) |
| CI | — | GitHub Actions | — |
| Deploy | — | RedHat OpenShift | oc CLI, Quay.io registry |

### Key files to read first when resuming

- `fintrago-project-guide.md` — full file hierarchy, tech stack, and prerequisite study lists
- `fintrago-product-doc.md` — this document
- `go-api/routes/routes.go` — all registered routes
- `ai-service/agents/graph.py` — LangGraph agent definition
- `docker-compose.yml` — how all services connect

