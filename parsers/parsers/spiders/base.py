import scrapy
import os

class Base(scrapy.Spider):
    PG_CONN_URI = os.environ['DB_CONN']
    HOSTNAME = 'https://chaturbate.com'
