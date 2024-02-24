FROM debian:stable-slim
COPY Barkin /bin/Barkin
ENV PORT 8000
CMD ["/bin/Barkin"]