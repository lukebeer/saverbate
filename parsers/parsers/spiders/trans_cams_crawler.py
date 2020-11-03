from .base_cams_crawler import BaseCamsCrawler

class TransCamsCrawler(BaseCamsCrawler):
    name = 'trans_cams_crawler'

    def gender(self):
      return 'trans'
