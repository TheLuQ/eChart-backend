from dataclasses import dataclass
from googleapiclient.discovery import build
import io
from googleapiclient.http import MediaIoBaseDownload


@dataclass
class DriveFile:
    id: str
    name: str

class DriveConnector:
    def __init__(self):
        self.service = build('drive', 'v3')

    def list_files(self, folder_id):
        query = f"'{folder_id}' in parents and trashed=false and mimeType = 'application/pdf'"
        results = self.service.files().list(q=query, fields="files(id, name)",).execute()
        return [DriveFile(**file) for file in results.get('files', [])]
    
    def get_file_info(self, file_id):
        file_info = self.service.files().get(fileId=file_id, fields="id, name").execute()
        return DriveFile(**file_info)

    def download_file(self, file_id):
        request = self.service.files().get_media(fileId=file_id)
        file_data = io.BytesIO()
        downloader = MediaIoBaseDownload(file_data, request)
        done = False
        while not done:
            status, done = downloader.next_chunk()
            print(f"Download {int(status.progress() * 100)}%.")
        file_data.seek(0)
        return file_data
