FROM scratch
MAINTAINER "oldthreefeng <louisehong4168@gmail.com>"
COPY stress /usr/bin/stress

ENTRYPOINT ["/usr/bin/stress"]
