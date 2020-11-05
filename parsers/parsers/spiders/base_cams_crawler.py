import scrapy
import postgresql
import re
from .base import Base

class BaseCamsCrawler(Base):
    SLASHES_RE = re.compile(r"^\/([^/]+)\/$")
    PAGE_RE = re.compile(r"^\/[^\/]+\/\?page=(\d+)$")
    """
    Maximum page to crawl
    """
    MAX_PAGE = 3

    name = 'female_cams_crawler'

    def start_requests(self):
        return [scrapy.Request(url=self.HOSTNAME+self.page_list_url(), callback=self.parse_cams)]

    def parse_cams(self, response):
        with postgresql.open(self.PG_CONN_URI) as db:
            for link in response.css('.details>.title'):
                path = link.css('a').attrib['href']
                username = self.SLASHES_RE.sub(r"\1", path)
                cam_fill = db.prepare(
                    "INSERT INTO broadcasters (name, created_at, gender) VALUES ($1, NOW(), $2) "
                    "ON CONFLICT ON CONSTRAINT uniq_broadcasters_name "
                    "DO NOTHING"
                )
                cam_fill(username, self.gender())

                yield {
                    'username': username
                }

        for next_page in response.css('a.endless_page_link'):
            href = next_page.attrib['href']

            if href != '/' and href != self.page_list_url():
                page_num = int(self.PAGE_RE.sub(r"\1", href))
                if page_num > self.MAX_PAGE:
                    continue
            else:
                continue

            yield response.follow(next_page, self.parse_cams)

    def page_list_url(self):
        return '/'+self.gender()+'-cams/'

    def gender(self):
      return 'female'
