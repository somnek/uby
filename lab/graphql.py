import requests
import json
from rich import print
from dotenv import load_dotenv
import os

load_dotenv()

owner = "charmbracelet"
repo = "bubbletea"
token = os.getenv("PERSONAL_ACCESS_TOKEN") or ""
url = "https://api.github.com/graphql"

# query to get all dependencies
query = """
query ($owner: String!, $repo: String!) {
  repository(owner: $owner, name: $repo) {
    dependencyGraphManifests(first: 100) {
      nodes {
        dependencies(first: 100) {
          nodes {
            packageName
            requirements
            repository {
              nameWithOwner
            }
          }
        }
      }
    }
  }
}
"""

variables = {"owner": owner, "repo": repo, "cursor": None}


headers = {
    "Authorization": "token " + token,
    "Accept": "application/vnd.github.hawkgirl-preview+json",
}
r = requests.post(url, json={"query": query, "variables": variables}, headers=headers)
data = json.loads(r.text)
print(data)
