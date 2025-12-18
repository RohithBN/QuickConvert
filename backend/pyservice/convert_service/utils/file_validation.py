from fastapi import UploadFile, HTTPException
from typing import List

def check_no_of_files(files: List[UploadFile], max_no: int):
    """
    Check the number of files uploaded and raise error accordingly
    max_no: passed from the API endpoint function
    """
    if len(files) > max_no:
        raise HTTPException(
            status_code=400, 
            detail=f"Maximum {max_no} files allowed."
        )
    
def check_file_types(
    files: List[UploadFile],
    allowed_types: set,
    allowed_exts: set
):
    """
    Checks if each file's content_type or extension is in the allowed sets.
    Raises HTTPException if any file is invalid.
    """
    for file in files:
        filename = file.filename.lower()
        ext = "." + filename.rsplit(".", 1)[-1] if "." in filename else ""
        if (file.content_type not in allowed_types) and (ext not in allowed_exts):
            raise HTTPException(
                status_code=400,
                detail=f"Invalid file type for '{file.filename}'. Allowed types: {allowed_types}, extensions: {allowed_exts}"
            )

def check_file_sizes(files: List[UploadFile], max_size_bytes: int):
    for file in files:
        file_size_bytes = file.size
        if file_size_bytes is not None and file_size_bytes > max_size_bytes:
            raise HTTPException(
                status_code=400,
                detail=f"File '{file.filename}' exceeds the maximum size of {max_size_bytes // (1024 * 1024)} MB"
            )

def validate_uploads(
    files: List[UploadFile],
    max_no: int,
    allowed_types: set,
    allowed_exts: set,
    max_size_bytes: int
):
    """
    Wrapper to validate uploaded files for count, type, and size.
    Raises HTTPException if any check fails.
    """
    check_no_of_files(files, max_no)
    check_file_types(files, allowed_types, allowed_exts)
    check_file_sizes(files, max_size_bytes)