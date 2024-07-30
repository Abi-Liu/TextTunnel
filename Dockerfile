FROM debian:stable-slim
COPY linux_binary /bin/server
CMD ["/bin/server"]
