#!/bin/sh

# Install ffmpeg on container start
apk add --no-cache ffmpeg

# Run the app
./app