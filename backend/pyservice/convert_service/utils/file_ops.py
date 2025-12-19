from fastapi import UploadFile
from pathlib import Path
from typing import List, Union
import uuid
import shutil
import zipfile

def save_uploaded_file(upload_file: UploadFile, destination: Path) -> Path:
    """
    Saves an UploadFile to the given destination with a unique name and
    returns the path to the saved file
    """
    short_id = str(uuid.uuid4())[:8]
    unique_name = f"{short_id}_{upload_file.filename}"
    file_path = destination/unique_name
    with open(file_path, "wb") as buffer:
        buffer.write(upload_file.file.read())

    return file_path

def zip_output_folders(output_folders: List[Path], work_dir: Path) -> Path:
    """
    Zips all folders in output_folders into a single zip file in work_dir and 
    returns the path to the zip file
    """
    zip_path = work_dir/"converted.zip"
    with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as zipf:
        for folder in output_folders:
            for file in folder.iterdir():
                arcname = f"{folder.name}/{file.name}"
                zipf.write(file, arcname=arcname)
                
    return zip_path

def cleanup_files(paths: List[Union[str, Path]]):
    """
    Deletes files or directories given in the paths list.
    Ignores errors if a file/folder does not exist.
    """
    for path in paths:
        p = Path(path)
        try:
            if p.is_file():
                p.unlink()
            elif p.is_dir():
                shutil.rmtree(p)
        except Exception as e:
            print(f"Cleanup failed for {p}: {e}")