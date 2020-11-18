#!/bin/bash

set -e

n=16
path_to_video=/app/downloads/$2/$1.mp4
duration=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 $path_to_video)
fps=$(echo "$n / $duration" | bc -l)
offset=$(echo "$duration / $n / 2" | bc -l)

echo "Duration: $duration, fps: $fps, offset: $offset"

ffmpeg -i $path_to_video -an -ss $offset -f image2 -vf fps=fps=$fps /app/downloads/$2/$1_thumb%04d.jpg
