FROM scratch

COPY stress /usr/bin/stress

ENTRYPOINT ["/usr/bin/stress"]
