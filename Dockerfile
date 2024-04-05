FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

ADD Verses /usr/bin/Verses

EXPOSE 80

CMD ["Verses"]


CMD ["Verses"]