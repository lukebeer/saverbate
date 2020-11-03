from .base_cams_crawler import BaseCamsCrawler

class CoupleCamsCrawler(BaseCamsCrawler):
    name = 'couple_cams_crawler'

    def gender(self):
      return 'couple'
