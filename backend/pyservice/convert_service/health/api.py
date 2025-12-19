from fastapi import APIRouter

router = APIRouter()
# no need of any prefix because its a single route - /health

@router.get("/health")
def health_check():
    return {"status": "ok"}