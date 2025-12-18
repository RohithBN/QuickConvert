from fastapi import UploadFile
from typing import List
from pathlib import Path
import subprocess

from convert_service.utils.file_ops import (
    save_uploaded_file,
    zip_output_folders
)

def convert_docxs_to_pdfs(
    files: List[UploadFile],
    work_dir: Path      # temp working directory
) -> tuple:
    """
    Converts each DOCX in `files` to PDF,
    stores each PDF in a separate folder inside `work_dir`
    and zips the folders and returns a path for download.
    Returns (zip_path, output_folders, docx_paths, errors)
    """
    errors = []

    # Save each file to disk (temp)
    docx_paths = []
    for upload_file in files:
        try:
            docx_path = save_uploaded_file(upload_file, work_dir)
            docx_paths.append(docx_path)
        except Exception as e:
            errors.append(f"Failed to save {upload_file.filename}: {str(e)}")

    # Convert each DOCX to PDF using LibreOffice
    output_folders = []
    for docx_path in docx_paths:
        try:
            # Create subfolder for current DOCX's PDF
            docx_stem = docx_path.stem
            output_folder = work_dir / docx_stem
            output_folder.mkdir(exist_ok=True)
            output_folders.append(output_folder)

            # Convert DOCX to PDF using LibreOffice
            result = subprocess.run([
                "soffice",
                "--headless",
                "--convert-to", "pdf",
                "--outdir", str(output_folder),
                str(docx_path)
            ], capture_output=True)
            if result.returncode != 0:
                raise RuntimeError(f"LibreOffice conversion failed: {result.stderr.decode()}")

        except Exception as e:
            errors.append(f"Failed to convert {docx_path.name}: {str(e)}")

    if not output_folders:
        raise RuntimeError("No DOCXs were successfully converted.")

    # Zip the entire thing and return it along with errors
    zip_path = zip_output_folders(output_folders, work_dir)
    return zip_path, output_folders, docx_paths, errors