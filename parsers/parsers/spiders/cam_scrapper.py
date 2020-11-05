import scrapy
import postgresql
import re
from .base import Base

class CamScrapper(Base):
    CURSOR_BATCH_SIZE = 50
    CURSOR_ID = 'cam_scrapper_cursor'
    name = 'cam_scrapper'

    def start_requests(self):
        with postgresql.open(self.PG_CONN_URI) as db:
            db.execute(
              "DECLARE " + self.CURSOR_ID + " CURSOR WITH HOLD FOR "
              "SELECT id, name FROM broadcasters WHERE "
              "language IS NULL OR "
              "followers IS NULL OR "
              "location IS NULL OR "
              "age IS NULL ORDER BY created_at DESC LIMIT 1000"
            )

            c = db.cursor_from_id(self.CURSOR_ID)
            urls = c.read(self.CURSOR_BATCH_SIZE)

            while len(urls) > 0:
                for item in urls:
                    path = "/"+item[1].strip()+"/"
                    doc_id = item[0]

                    url = self.HOSTNAME + path
                    yield scrapy.Request(url=url, meta={'doc_id': doc_id, 'name': item[1]})
                urls = c.read(self.CURSOR_BATCH_SIZE)
            c.close()

    def parse(self, response):
        doc_id = response.meta['doc_id']
        username = response.meta['name']

        followers = ''.join(
          response.xpath(self.xpath_attribute('Followers')).getall()
        )
        language = ''.join(
          response.xpath(self.xpath_attribute('Languages')).getall()
        )
        location = ''.join(
          response.xpath(self.xpath_attribute('Location')).getall()
        )
        age = ''.join(
          response.xpath(self.xpath_attribute('Age')).getall()
        )
        if followers == '':
          followers = '0'
        if age == '':
          age = '0'

        with postgresql.open(self.PG_CONN_URI) as db:
            try:
              broadcaster_fill = db.prepare(
                "UPDATE broadcasters SET "
                "language = $1, "
                "followers = $2::int4, "
                "location = $3, "
                "age = $4::int4 "
                "WHERE id = $5"
              )
              broadcaster_fill(language, int(followers), location, int(age), doc_id)
            except:
              print({
                'doc_id': doc_id,
                'username': username,
                'language': language,
                'followers': followers,
                'location': location,
                'age': age,
              })
              raise

        yield {
          'doc_id': doc_id,
          'username': username,
          'language': language,
          'followers': followers,
          'location': location,
          'age': age,
        }

    def xpath_attribute(self, attribute):
      return '//div[@class="attribute"]//div[@class="label" and contains(text(),"'+attribute+':")]/following-sibling::div[@class="data"]/text()'
