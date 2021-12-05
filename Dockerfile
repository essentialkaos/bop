## BUILDER #####################################################################

FROM essentialkaos/centos:7 as builder

RUN mkdir -p "/go/src" && chmod -R 777 "/go"

ENV GOPATH=/go

WORKDIR /go/src/github.com/essentialkaos/bop

COPY . .

# hadolint ignore=DL3031,DL3033
RUN yum -y -q install https://yum.kaos.st/kaos-repo-latest.el7.noarch.rpm && \
    yum -y update && \
    yum -y -q install make golang git upx && \
    yum clean all && \
    make deps && \
    make all && \
    upx bop

## FINAL IMAGE #################################################################

FROM essentialkaos/centos:7

LABEL name="Bop Image on CentOS 7" \
      vendor="ESSENTIAL KAOS" \
      maintainer="Anton Novojilov" \
      license="Apache-2.0" \
      version="2021.07.21"

COPY --from=builder /go/src/github.com/essentialkaos/bop/bop /usr/bin/

VOLUME /bop
WORKDIR /bop

ENTRYPOINT ["bop"]

################################################################################
