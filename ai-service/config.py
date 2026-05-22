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
