FROM docker.io/library/golang:1.17 AS builder

WORKDIR /builder
COPY ./src/main.go /builder/
RUN CGO_ENABLED=0 go build -o wellknown -ldflags '-extldflags "-static" -w -s'  main.go

FROM scratch

COPY --from=builder /builder/wellknown /bin/wellknown
EXPOSE 80
USER 1000

ENTRYPOINT [ "/bin/wellknown"]