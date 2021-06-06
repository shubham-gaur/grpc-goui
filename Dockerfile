FROM  alpine:3
LABEL authors="Shubham Gaur <shubham-gaur.github.io>"
WORKDIR /root/goui
ENV PATH "$PATH:/root/goui/assets/" 
ENV EXPOSED_PORTS=8081...8090 \
    INTERFACE=eth0
COPY assets ./assets
COPY templates ./templates
COPY bin/goui ./
ENTRYPOINT  ["./goui"]