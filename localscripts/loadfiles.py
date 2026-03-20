import re
import unicodedata

from drive import DriveConnector, DriveFile
from firestore import FireConnector
from storage import StorageConnector
from parserfinal import createInstrument

driveConnector = DriveConnector()


def _upload_to_storage(file: DriveFile, storageConnector: StorageConnector, titlemain: str, path: str | None = None):
    instrument = createInstrument(file.name, file.id)
    if not instrument:
        print(f"Failed to create instrument for file '{file.name}'")
        return
    if path is None:
        path = re.sub(r'\s+', '-', titlemain).lower()
    path = unicodedata.normalize('NFKD', path).encode('ascii', 'ignore').decode('ascii')
    file_data = driveConnector.download_file(file.id)
    storageConnector.upload_file(
        file_data, path,
        instrument.instrument_name_pol, instrument.instrument_name_en,
        instrument.voice, instrument.key, instrument.file_id,
    )
    print(f"File '{file.name}' uploaded to storage with metadata: {instrument.__dict__}")


def load_drive_file(file_id: str, bucket_name: str, titlemain: str, custom_path: str | None = None):
    storageConnector = StorageConnector(bucket_name)
    file = driveConnector.get_file_info(file_id)
    _upload_to_storage(file, storageConnector, titlemain, custom_path)


def load_from_drive_folder(folder_id: str, collection_name: str, bucket_name: str, custom_title: str | None = None):
    files = driveConnector.list_files(folder_id)
    title = driveConnector.get_file_info(folder_id).name
    polish_chars = str.maketrans("ąćęłńóśźżĄĆĘŁŃÓŚŹŻ", "acelnoszzACELNOSZZ")
    gcp_folder = "-".join(part for part in title.translate(polish_chars).lower().split() if part)
    print(f"Found {len(files)} PDF files in folder '{folder_id}'")
    print(f"Title of the folder: {title}")
    storageConnector = StorageConnector(bucket_name)
    for file in files:
        _upload_to_storage(file, storageConnector, title, gcp_folder)
    print("Loading to storage completed.")
    with FireConnector(collection_name) as fireConnector:
        fireConnector.set_title(gcp_folder, custom_title or title)
