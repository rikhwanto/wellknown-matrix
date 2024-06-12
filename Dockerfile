FROM docker.io/library/golang:1.22 AS builder

WORKDIR /builder
COPY ./src/ /builder/
RUN CGO_ENABLED=0 go build -o wellknown -ldflags '-extldflags "-static" -w -s'  ./...

FROM scratch

COPY --from=builder /builder/wellknown /bin/wellknown
EXPOSE 8080
USER 1000

ENTRYPOINT [ "/bin/wellknown"]