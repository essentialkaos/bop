## BUILDER #####################################################################

FROM centos:7 as builder

RUN mkdir -p "/go/src" && chmod -R 777 "/go"

ENV GOPATH=/go

WORKDIR /go/src/github.com/essentialkaos/bop

COPY . .

RUN yum -y -q install https://yum.kaos.st/kaos-repo-latest.el7.noarch.rpm && \
    yum -y -q install make golang git upx && \
    make deps && \
    make all && \
    upx bop

## FINAL IMAGE #################################################################

FROM centos:7

LABEL name="Bop Image on CentOS 7" \
      vendor="ESSENTIAL KAOS" \
      maintainer="Anton Novojilov" \
      license="EKOL" \
      version="2020.05.09"

COPY --from=builder /go/src/github.com/essentialkaos/bop/bop /usr/bin/
COPY --from=builder /go/src/github.com/essentialkaos/bop/bop-entrypoint /usr/bin/

VOLUME /bop
WORKDIR /bop

ENTRYPOINT ["bop"]

################################################################################
