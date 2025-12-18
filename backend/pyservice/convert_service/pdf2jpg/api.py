from fastapi import APIRouter, UploadFile, File, HTTPException, BackgroundTasks
from fastapi.responses import FileResponse
from typing import List
from pathlib import Path

from convert_service.utils.file_validation import validate_uploads
from convert_service.pdf2jpg.service import convert_pdfs_to_jpgs
from convert_service.utils.file_ops import cleanup_files

router = APIRouter(
    prefix="/pdf-to-jpg",
    tags=["pdf", "jpg"]
)

MAX_FILES = 50
MAX_FILE_SIZE_MB = 20
MAX_FILE_SIZE_BYTES = MAX_FILE_SIZE_MB * 1024 * 1024

@router.post("/", summary="PDFs to JPGs conversion")
async def convert_pdf_to_jpg(files: List[UploadFile] = File(...), background_tasks: BackgroundTasks = None):
    validate_uploads(
        files,
        max_no = MAX_FILES,
        allowed_types = {"application/pdf"},
        allowed_exts = {".pdf"},
        max_size_bytes = MAX_FILE_SIZE_BYTES
    )

    # Use directory `temp_processing` and call the service function
    work_dir = Path(__file__).parent.parent.parent/"temp_processing"
    work_dir.mkdir(exist_ok=True)
    try:
        zip_path, output_folders, pdf_paths, errors = convert_pdfs_to_jpgs(files, work_dir)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Conversion failed: {str(e)}")

    # Schedule cleanup of all temp files/folders after response
    cleanup_paths = [zip_path] + output_folders + pdf_paths
    if background_tasks is not None:
        background_tasks.add_task(cleanup_files, cleanup_paths)

    response = FileResponse(
        path=zip_path,
        filename="converted_images.zip",
        media_type="application/zip"
    )
    if errors:
        # Only include short error messages in header
        response.headers["X-Conversion-Errors"] = "; ".join(errors)
    return response