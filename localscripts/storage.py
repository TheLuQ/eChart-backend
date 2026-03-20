from google.cloud import storage
from google.cloud.storage import Bucket

class StorageConnector:
    def __init__(self, bucket_name: str):
        self.connector = storage.Client()
        self.bucket = self.connector.bucket(bucket_name)

    def upload_file(self, file, destination_path: str, instrument_pl: str,
                instrument_en: str, voice: str | None, key: str | None, file_id: str
                ):
        file_name = '-'.join(filter(None, [instrument_en, voice, key])) + '.pdf'
        blob = self.bucket.blob(f'{destination_path}/{file_name}')
        blob.metadata = {
            "instrument_name_pol": instrument_pl,
            "instrument_name_en": instrument_en,
            **({"voice": voice} if voice else {}),
            **({"key": key} if key else {}),
            "file_id": file_id
        }
        blob.content_type = 'application/pdf'
        blob.upload_from_file(file)
