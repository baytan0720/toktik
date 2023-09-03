FROM alpine:latest

WORKDIR /
COPY main .
COPY etc/config.yaml.sample /etc/config.yaml
COPY ffmpeg /bin/ffmpeg
CMD ["./main", "-config", "/etc/config.yaml", "-consul","consul:8500"]