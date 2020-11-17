#!/bin/bash

set -e

ffmpeg -i /app/downloads/$2/$1.mp4 -an -ss 00:00:10.000 -vframes 1 /app/downloads/$2/$1.jpg
