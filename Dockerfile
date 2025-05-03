FROM --platform=linux/amd64 grafana/grafana-enterprise:11.6.1

COPY ./dist /var/lib/grafana/plugins/dalelane-kafka-datasource
