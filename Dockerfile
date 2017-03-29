FROM alpine:latest

COPY ./bin/go-deploy /usr/local/bin/go-deploy/

COPY ./run.sh /

ENV PATH="$PATH:/usr/local/bin/go-deploy"

RUN chmod +x /usr/local/bin/go-deploy/go-deploy

CMD ["/bin/sh", "/run.sh"]