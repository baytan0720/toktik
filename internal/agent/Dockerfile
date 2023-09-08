FROM minio/minio:latest

ARG TZ=Asia/Shanghai
ARG MINIO_ACCESS_KEY=admin
ARG MINIO_SECRET_KEY=Aa112211

ENV TZ=${TZ}
ENV MINIO_ROOT_USER=${MINIO_ACCESS_KEY}
ENV MINIO_ROOT_PASSWORD=${MINIO_SECRET_KEY}

COPY main .
COPY etc/config.yaml.sample /etc/config.yaml
COPY ffmpeg /bin/ffmpeg
ENTRYPOINT ["./main", "-config", "/etc/config.yaml"]