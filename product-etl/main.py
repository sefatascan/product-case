import json

import requests

if __name__ == '__main__':

    with open('sample.json', 'r') as file:
        samples = json.load(file)

        bulk_data = ""
        for sample in samples:
            bulk_data += json.dumps({"index": {"_index": "product-index"}}) + "\n"
            bulk_data += json.dumps(sample) + "\n"

        response = requests.post(f"http://localhost:9200/product-index/_bulk",
                                 headers={"Content-Type": "application/json"},
                                 data=bulk_data)

        if not response.ok:
            raise Exception("Elasticsearch documents could not be inserted")
