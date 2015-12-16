FROM gliderlabs/alpine

ADD build/csv2api.linux /

ENV PORT 8080
ENV SERVE_FROM /tmp/data
ENV API_KEY ""

VOLUME ["/tmp/data"]
EXPOSE 8080/tcp

CMD ["/csv2api.linux"]