FROM nginx:1.29.0-alpine3.22

RUN apk add --no-cache inotify-tools

COPY ./watch-certs.sh /watch-certs.sh
RUN chmod +x /watch-certs.sh

CMD ["/bin/sh", "-c", "nginx && /watch-certs.sh"]
