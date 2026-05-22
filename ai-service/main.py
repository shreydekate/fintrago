from fastapi import FastAPI
from routers import upload, ask

app = FastAPI(title="Fintrago AI Service")

app.include_router(upload.router, prefix="/upload", tags=["upload"])
app.include_router(ask.router, prefix="/ask", tags=["ask"])


@app.get("/health")
def health():
    return {"status": "ok"}
