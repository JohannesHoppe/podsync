# multi stage dockerfile
#
# build this file with the following command
# docker buildx build -t podsync-test --platform=linux/amd64 .
#
# Hint, this does not work for me:
# docker buildx build -t podsync-test --platform=linux/amd64,linux/arm64,linux/arm/v6,linux/386 .
# > error: multiple platforms feature is currently not supported for docker driver. Please switch to a different driver (eg. "docker buildx create --use")


# building the go binary
# see https://github.com/mxpv/podsync/issues/56#issuecomment-717777668
FROM golang:alpine AS builder
LABEL stage=builder
WORKDIR /workspace
COPY . .
RUN go build -o /bin/podsync ./cmd/podsync


# The actual podsync Dockerfile, but with yt-dlp instead of youtube-dl
# see https://github.com/tuxpeople/docker-podsync/blob/a27674c692fe9a27dd43ef27685ad6440dbd8726/Dockerfile

FROM alpine:3.16.0
WORKDIR /app/
# hadolint ignore=DL3018,DL3017
RUN apk --no-cache upgrade && \
    apk --no-cache add ca-certificates ffmpeg tzdata python3 && \
    wget -q -O /usr/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp && \
    chmod +x /usr/bin/yt-dlp && \
    ln -s /usr/bin/yt-dlp /usr/bin/youtube-dl
COPY --from=builder /bin/podsync .
CMD ["/app/podsync"]


