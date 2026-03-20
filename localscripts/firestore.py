from google.cloud.firestore import Client

class FireConnector:
    def __init__(self, collection_name):
        self.collection_name = collection_name
        self.db = Client()

    def load_data_col(self, data, doc_ref=None):
        print(f"Loading {len(data)} documents into Firestore collection '{self.collection_name}'...")
        for doc in data:
            self.bulk_writer.set(self.db.collection(self.collection_name).document(doc_ref), doc)
        self.bulk_writer.flush()

    def load_data(self, data, doc_ref=None):
        print(f"Loading documents into Firestore collection '{self.collection_name}'...")
        self.bulk_writer.set(self.db.collection(self.collection_name).document(doc_ref), data)
        self.bulk_writer.flush()

    def search_by_parent_path(self, path: str):
        docs = self.db.collection(self.collection_name).where('parent_path', '==', path).get()
        return [doc.reference for doc in docs]

    def set_title(self, path: str, title: str):
        for ref in self.search_by_parent_path(path):
            ref.set({"title": title}, merge=True)

    def __enter__(self):
        self.bulk_writer = self.db.bulk_writer()
        return self
    
    def __exit__(self, exc_type, exc, tb):
        print("Flushing remaining writes...")
        self.bulk_writer.close()
        pass
