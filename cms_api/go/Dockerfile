FROM golang:1.24.1-bullseye as build-env

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/api

FROM alpine:latest

COPY --from=build-env /go/bin/app /

# ここが重要
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.6.0 /lambda-adapter /opt/extensions/lambda-adapter

# ここが重要
ENV PORT=3000
EXPOSE 3000

CMD ["/app"]
