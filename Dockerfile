FROM scratch
MAINTAINER Eric Izoita <ericizoita@gmail.com>
EXPOSE 8080
COPY resources/ /resources/
COPY coalesce.cfg /
COPY bin/coalesce /
ENTRYPOINT ["/coalesce"]
