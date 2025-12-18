from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from convert_service.health.api import router as health_router
from convert_service.pdf2jpg.api import router as pdf2jpg_router
from convert_service.docx2pdf.api import router as docx2pdf_router

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["POST", "GET", "OPTIONS"],
    allow_headers=["Origin", "Content-Type"],
    expose_headers=["Content-Length"],
)

app.include_router(health_router)
app.include_router(pdf2jpg_router)
app.include_router(docx2pdf_router)