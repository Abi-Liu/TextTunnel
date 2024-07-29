FROM debian:stable-slim
COPY linux_binary /bin/server
COPY .env ./
CMD ["/bin/server"]
