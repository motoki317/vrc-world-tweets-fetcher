FROM golang:1.18-alpine AS build

WORKDIR /work
ENV CGO_ENABLED 0
ENV GOCACHE=/tmp/go/cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/tmp/go/cache go build -o main -ldflags "-s -w" .

FROM alpine:latest AS runner

WORKDIR /work

RUN apk --no-cache add curl bash
RUN curl -o wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh
RUN chmod 0744 ./wait-for-it.sh
COPY --from=build /work/main ./

ENTRYPOINT ["./main"]
