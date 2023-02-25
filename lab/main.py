import requests
import json
from rich import print
from dotenv import load_dotenv
import os

load_dotenv()

owner = 'charmbracelet'
repo = 'bubbletea'
token = os.getenv('PERSONAL_ACCESS_TOKEN') or ""
url = 'https://api.github.com/graphql'

query = '''
query ($owner: String!, $repo: String!) {
    repository(owner: "charmbracelet", name: "bubbletea") {
      dependents(first: 100, after: null) {
        node {
          nameWithOwner
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
    }
  }
'''


headers = {'Authorization': 'token ' + token}
r = requests.post(url, json={'query': query}, headers=headers)
data = json.loads(r.text)
print(data)


