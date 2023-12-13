FROM golang:1.21-alpine AS builder
LABEL authors="antony"
RUN apk add --no-cache make
WORKDIR /usr/src/app
COPY . .
RUN make dep-download; make

FROM golang:1.21-alpine
COPY --from=builder /usr/src/app/analytics /usr/bin/analytics
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/health || exit 1
EXPOSE 8080
ENTRYPOINT [ "/usr/bin/analytics" ]