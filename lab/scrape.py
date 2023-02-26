import requests
from typing import List, Text, Tuple, Dict
from bs4 import BeautifulSoup
from requests.adapters import HTTPAdapter
from requests.packages.urllib3.util.retry import Retry
from rich import console

console = console.Console()
PACKAGE_ID = 'UGFja2FnZS0yMjc1ODk0MDQy'
PACKAGE_NAME = 'github.com/charmbracelet/bubbletea'
URL = f'https://{PACKAGE_NAME}/network/dependents?package_id={PACKAGE_ID}'

package = {'id': PACKAGE_ID, 'name': PACKAGE_NAME, 'url': URL}

class Scraper:
  def __init__(self, placeholder) -> None:
    self.placeholder = placeholder

  def collect(self):
    next_exists = True
    result = []
    page_number = 0

    # Get dependent count
    r = self.requests_retry_session().get(URL)
    soup = BeautifulSoup(r.content, 'html.parser')
    svg_item = soup.find('svg', {'class': 'octicon-code-square'})
    a_around_svg = svg_item.parent
    total_dependencies = self.to_int(a_around_svg.text.strip().split()[0])

    # Parse all dependent packages pages
    while next_exists:
      r = self.requests_retry_session().get(URL)
      soup = BeautifulSoup(r.content, 'html.parser')

      # Browse page dependents
      for t in soup.find_all('div', {'class': 'Box-row'}):
        owner = t.find('a', {'data-repository-hovercards-enabled': ''}).text
        repo = t.find('a', {'data-hovercard-type': 'repository'}).text
        stars = self.to_int(t.find('svg', {'class': 'octicon-star'}).parent.text.strip())
        forks = self.to_int(t.find('svg', {'class': 'octicon-repo-forked'}).parent.text.strip())
        
        result_item = {
          'name': f'{owner}/{repo}',
          'stars': stars,
          'forks': forks,
        }
        result += [result_item]

      # Check next page
      paginate_container = soup.find('div', 'paginate-container')
      if paginate_container:
        for u in paginate_container.find_all('a'):
          if u.text == 'Next':
            next_exists = True
            page_number += 1
            console.log(f'📑: {page_number}')
      
      break
    
  @staticmethod
  def to_int(s: str) -> int:
    return ''.join(c for c in s if c.isdigit())

  def requests_retry_session(
    self,
    retries: int = 3,
    backoff_factor: float = 0.5,
    status_forcelist: List[Tuple] = (500, 502, 504),
    session: requests.Session = None,
  ) -> requests.Session:

    session = session or requests.Session()
    retry = Retry(
      total=retries,
      read=retries,
      connect=retries,
      backoff_factor=backoff_factor,
      status_forcelist=status_forcelist
    )
    adapter = HTTPAdapter(max_retries=retry)
    session.mount('http://', adapter)
    session.mount('https://', adapter)
    return session

if __name__ == '__main__':
  scrape = Scraper(placeholder="x")
  scrape.collect()
