# Saverbate

Saverbate is complex service for parsing, scrapping, and recording translations from chaturbate.com.

Services:

1) crawler - service for crawling performers and add them and their data to db
2) subscriber - service for subscribing to performers
3) mailer - service for checking who is online and sending signal to event bus about it
4) downloader - service for downloading performer translation
5) frontend - web server

## Subscriber

On development stage. Main issue is Cloudflare hCaptcha appearanced after 4th similar action on the site.

## Mailer

Now successfull checks who from performers is online. TODO: move read envelope to bin.

## Downloader ⬇

The core of the system. Work in progress.
It should receive an event about the online performer, start streamlink and download the stream until the performer disconnect.
Then the completed video file must be passed to ffmpeg for conversion to hls format.
Converted video should be uploaded to aws s3 server.

## Frontend 🐼

Part of the core of the system. Shows hls videos on website.
Work in progress. Only the web server base is ready.
