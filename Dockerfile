FROM scratch
MAINTAINER Eric Izoita <ericizoita@gmail.com>
EXPOSE 8080
COPY resources/ /resources/
COPY bin/coalesce /
#ENV ADMIN_PASS password
ENV SESSION_SECRET correcthorsebatterystaple
ENTRYPOINT ["/coalesce"]
