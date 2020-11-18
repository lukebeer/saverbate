FROM jrottenberg/ffmpeg:4.2-ubuntu AS ffmpeg

FROM python:3.8.6-buster AS py

COPY --from=ffmpeg /usr/local /usr/local/
ENV LD_LIBRARY_PATH=/usr/local/lib:/usr/local/lib64

RUN apt-get update -yq && \
  apt-get install -yq \
    tini \
    curl \
    git \
    bc

ENTRYPOINT ["tini", "--"]

RUN pip install youtube-dl
