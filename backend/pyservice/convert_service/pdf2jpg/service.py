from fastapi import UploadFile
from typing import List
from pathlib import Path

from convert_service.utils.file_ops import (
    save_uploaded_file,
    zip_output_folders
)
from pdf2image import convert_from_path

def convert_pdfs_to_jpgs(
    files: List[UploadFile],
    work_dir: Path      # temp working directory
) -> tuple:
    """
    Converts each PDF in `files` to JPG images page-by-page,
    stores each PDF's images in a separate folder inside `work_dir`
    and zips the folders and returns a path for download
    """
    errors = []

    # Save each file to disk (temp)
    pdf_paths = []
    for upload_file in files:
        try:
            pdf_path = save_uploaded_file(upload_file, work_dir)
            pdf_paths.append(pdf_path)
        except Exception as e:
            errors.append(f"Failed to save {upload_file.filename}: {str(e)}")

    # Convert each PDF to JPG using pdf2image
    output_folders = []
    for pdf_path in pdf_paths:
        try:
            # Create subfolder for current PDF images
            pdf_stem = pdf_path.stem
            output_folder = work_dir/pdf_stem
            output_folder.mkdir(exist_ok=True)
            output_folders.append(output_folder)

            # Convert PDF pages
            images = convert_from_path(str(pdf_path))
            for i, image in enumerate(images, start = 1):
                image_path = output_folder/f"page_{i}.jpg"
                image.save(image_path, "JPEG")
        except Exception as e:
            errors.append(f"Failed to convert {pdf_path.name}: {str(e)}")

    if not output_folders:
        raise RuntimeError("No PDFs were successfully converted.")
    
    # Zip the entire thing and return it along with errors
    zip_path = zip_output_folders(output_folders, work_dir)
    return zip_path, output_folders, pdf_paths, errors
