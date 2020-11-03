FROM jrottenberg/ffmpeg:4.2-ubuntu AS ffmpeg

FROM python:3.8.6-buster AS py

COPY --from=ffmpeg /usr/local /usr/local/
ENV LD_LIBRARY_PATH=/usr/local/lib:/usr/local/lib64

RUN apt-get update -yq && \
  apt-get install -yq \
    tini \
    curl \
    git

ENTRYPOINT ["tini", "--"]

RUN cd /tmp && \
  git clone https://github.com/isqad/streamlink && \
  cd /tmp/streamlink && \
  python setup.py install && \
  streamlink --version && \
  ffmpeg -version
