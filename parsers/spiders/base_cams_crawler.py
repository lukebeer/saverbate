import scrapy
import postgresql
import re
from .base import Base

class BaseCamsCrawler(Base):
    SLASHES_RE = re.compile(r"^\/([^/]+)\/$")

    name = 'female_cams_crawler'

    def start_requests(self):
        return [scrapy.Request(url=self.HOSTNAME+'/'+self.gender()+'-cams/', callback=self.parse_cams)]

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
            yield response.follow(next_page, self.parse_cams)

    def gender(self):
      return 'female'
