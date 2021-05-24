FROM golang:alpine3.12 as builder

ARG VCS_REF
ARG DB_URI_PATH=file:indahaus.db?_busy_timeout=5000

WORKDIR /golang
COPY . .

RUN apk add --no-cache sqlite gcc musl-dev

WORKDIR /golang/cmd/api
RUN go build -ldflags "-X main.build=${VCS_REF}" -o indahaus

# TODO: come up with a better way to generate db file
WORKDIR /golang/cmd
RUN DB_URI="${DB_URI_PATH}" go run admin/main.go migrate

FROM alpine:3.12

ARG BUILD_DATE
ARG VCS_REF
ARG PORT

RUN addgroup -g 1000 golang && \
  adduser -u 1000 -G golang -s /bin/sh -D golang

RUN mkdir /golang && chown golang:golang /golang
WORKDIR /golang

COPY --from=builder /usr/bin/sqlite3 /usr/bin/sqlite3
COPY --from=builder --chown=golang:golang /golang/cmd/indahaus.db ./
COPY --from=builder --chown=golang:golang /golang/cmd/api/indahaus ./
COPY config.yaml ./

ENV PORT="${PORT}"

CMD ["./indahaus"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
  org.opencontainers.image.title="indahaus" \
  org.opencontainers.image.authors="Shane Unger <shane.unger1@gmail.com>" \
  org.opencontainers.image.revision="${VCS_REF}"

USER golang