FROM scratch
COPY build/linux/restapi /
COPY build/linux/dbstore /
COPY config_linux.toml /config.toml
VOLUME /db
