FROM phusion/baseimage

COPY ./bin/watcher /watcher
COPY ./watcherd.docker.conf /watcherd.docker.conf

CMD ["./watcher", "--config=./watcherd.docker.conf"]