FROM python:3.8.6-buster

RUN apt-get update -yq && \
  apt-get install -yq \
    tini

ENTRYPOINT ["tini", "--"]

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
