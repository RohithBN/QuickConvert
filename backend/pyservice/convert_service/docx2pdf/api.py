from fastapi import APIRouter, UploadFile, File, HTTPException, BackgroundTasks
from fastapi.responses import FileResponse
from typing import List
from pathlib import Path

from convert_service.utils.file_validation import validate_uploads
from convert_service.docx2pdf.service import convert_docxs_to_pdfs
from convert_service.utils.file_ops import cleanup_files

router = APIRouter(
    prefix="/docx-to-pdf",
    tags=["docx", "pdf"]
)

MAX_FILES = 50
MAX_FILE_SIZE_MB = 20
MAX_FILE_SIZE_BYTES = MAX_FILE_SIZE_MB * 1024 * 1024

@router.post("/", summary="DOCXs to PDFs conversion")
def convert_docx_to_pdf(files: List[UploadFile] = File(...), background_tasks: BackgroundTasks = None):
    validate_uploads(
        files,
        max_no = MAX_FILES,
        allowed_types = {
            "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
            "application/msword"
        },
        allowed_exts = {".docx", ".doc"},
        max_size_bytes = MAX_FILE_SIZE_BYTES
    )
    
    # Use directory `temp_processing` and call the service function
    work_dir = Path(__file__).parent.parent/"temp_processing"
    work_dir.mkdir(exist_ok=True)
    try:
        zip_path, output_folders, docx_paths, errors = convert_docxs_to_pdfs(files, work_dir)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Conversion failed: {str(e)}")
    
    # Schedule cleanup of all temp files/folders after response
    cleanup_paths = [zip_path] + output_folders + docx_paths
    if background_tasks is not None:
        background_tasks.add_task(cleanup_files, cleanup_paths)

    response = FileResponse(
        path=zip_path,
        filename="converted_docs.zip",
        media_type="application/zip"
    )
    if errors:
        # Only include short error messages in header
        response.headers["X-Conversion-Errors"] = "; ".join(errors)
    return response