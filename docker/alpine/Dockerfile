FROM snowdreamtech/golang:1.23.5-alpine3.21 AS builder

# Switch to the user
USER root

# Set the workdir
WORKDIR /root

ENV ProjectName="gserver" \
    Author="Snowdream Tech <snowdreamtech@qq.com>" \
    COPYRIGHT="Copyright (c) 2023-present SnowdreamTech Inc." \
    LICENSE="MIT <https://github.com/snowdreamtech/gserver/blob/main/LICENSE>" \
    CmdPath="github.com/snowdreamtech/gserver/pkg/env"  \
    BuildTime=N/A \
    CommitHash=N/A \
    CommitHashFull=N/A \
    GoVersion=N/A  \
    GitTag=N/A \
    Debug=false

COPY . .

RUN apk add --no-cache \
    git

RUN go mod tidy \
    && BuildTime=$(date -u +'%Y.%m.%d %H:%M:%S %z') \
    && GoVersion=$(go version | awk '{print $3,$4}') \
    && CommitHash=$(git log -1 --pretty=format:%h 2>/dev/null || echo 'N/A') \
    && CommitHashFull=$(git log -1 --pretty=format:%H 2>/dev/null || echo 'N/A') \
    && GitTag=$(git describe --tags --abbrev=0 2>/dev/null|| echo 'N/A') \
    && LDGOFLAGS="-X '${CmdPath}.ProjectName=${ProjectName}'" \
    && LDGOFLAGS="${LDGOFLAGS} -X '${CmdPath}.Author=${Author}'" \
    && LDGOFLAGS="${LDGOFLAGS} -X '${CmdPath}.BuildTime=${BuildTime}'" \
    && LDGOFLAGS="${LDGOFLAGS} -X '${CmdPath}.CommitHash=${CommitHash}'" \
    && LDGOFLAGS="${LDGOFLAGS} -X '${CmdPath}.CommitHashFull=${CommitHashFull}'" \
    && LDGOFLAGS="${LDGOFLAGS} -X '${CmdPath}.GoVersion=${GoVersion}'" \
    && LDGOFLAGS="${LDGOFLAGS} -X '${CmdPath}.GitTag=${GitTag}'" \
    && LDGOFLAGS="${LDGOFLAGS} -X '${CmdPath}.COPYRIGHT=${COPYRIGHT}'" \
    && LDGOFLAGS="${LDGOFLAGS} -X '${CmdPath}.LICENSE=${LICENSE}'" \
    && go build -trimpath -ldflags "${LDGOFLAGS} -s -w" -v -o gserver 





FROM snowdreamtech/alpine:3.21.0

# OCI annotations to image
LABEL org.opencontainers.image.authors="Snowdream Tech" \
    org.opencontainers.image.title="GServer Image Based On Alpine" \
    org.opencontainers.image.description="Docker Images for GServer on Alpine. (i386, amd64, arm32v6, arm32v7, arm64, ppc64le,riscv64, s390x)" \
    org.opencontainers.image.documentation="https://hub.docker.com/r/snowdreamtech/gserver" \
    org.opencontainers.image.base.name="snowdreamtech/gserver:alpine" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.source="https://github.com/snowdreamtech/gserver" \
    org.opencontainers.image.vendor="Snowdream Tech" \
    org.opencontainers.image.version="0.0.5" \
    org.opencontainers.image.url="https://github.com/snowdreamtech/gserver"

# Switch to the user
USER root

# Set the workdir
WORKDIR /root

# keep the docker container running
ENV KEEPALIVE=1 \
    # The cap_net_bind_service capability in Linux allows a process to bind a socket to Internet domain privileged ports, 
    # which are port numbers less than 1024. 
    CAP_NET_BIND_SERVICE=0 

ARG GID=1000 \
    UID=1000  \
    USER=gserver \
    WORKDIR=/var/lib/gserver

# Create a user with UID and GID
RUN if [ "${USER}" != "root" ]; then \
    addgroup -g ${GID} ${USER}; \
    adduser -H -u ${UID} -g ${USER} -G ${USER} -s /sbin/nologin -D ${USER}; \
    # sed -i "/%sudo/c ${USER} ALL=(ALL:ALL) NOPASSWD:ALL" /etc/sudoers; \
    fi

# Enable CAP_NET_BIND_SERVICE
RUN if [ "${USER}" != "root" ] && [ "${CAP_NET_BIND_SERVICE}" -eq 1 ]; then \
    apk add --no-cache libcap; \
    # setcap 'cap_net_bind_service=+ep' `which nginx`; \
    fi

COPY --from=builder /root/gserver /usr/local/bin

RUN apk add --no-cache \
    vim \
    && mkdir -p ${WORKDIR} \
    && mkdir -p /var/log/gserver \
    && find / -name "*gserver*" -exec chown -Rv ${USER}:${USER} {} \; 

# Switch to the user
USER ${USER}

# Set the workdir
WORKDIR ${WORKDIR}

COPY docker/alpine/docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]

CMD [ "GIN_MODE=release","/usr/local/bin/gserver", "--wwwroot", "/var/lib/gserver", "--log-dir", "/var/log/gserver"]