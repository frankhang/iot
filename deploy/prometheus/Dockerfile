FROM prom/prometheus AS builder

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /etc/prometheus/console_libraries ./console_libraries
COPY --from=builder /etc/prometheus/consoles ./consoles
COPY --from=builder /bin/prometheus .
RUN mkdir data

ARG config=prometheus
COPY ${config}.* ./
RUN mv ${config}.yml prometheus.yml
# COPY ${config}.yml ./prometheus.yml
CMD ["./prometheus", "--config.file=prometheus.yml"]  