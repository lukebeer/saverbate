import scrapy

class Base(scrapy.Spider):
    PG_CONN_URI = 'pq://postgres:qwerty@localhost:10532/chaturbate_records'
    HOSTNAME = 'https://chaturbate.com'
