FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY .. ./

RUN go mod download


RUN CGO_ENABLED=0 go build -o /mentat-backend

FROM gcr.io/distroless/base-debian10 AS runner

WORKDIR /

COPY --from=builder /mentat-backend /mentat-backend

EXPOSE 8090

USER nonroot:nonroot

ENTRYPOINT ["/mentat-backend"]