FROM golang:1.23-alpine AS build
RUN apk update && \
    apk add curl \
    git \
    bash \
    make \
    ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.* ./
RUN go mod download
RUN go mod verify

# copy source files and build the binary
COPY . .
RUN make build


FROM alpine:latest
EXPOSE 3000
RUN apk --no-cache add ca-certificates bash
RUN mkdir -p /var/log/app
WORKDIR /app/
#COPY --from=build /usr/local/bin/migrate /usr/local/bin
#COPY --from=build /app/migrations ./migrations/
COPY --from=build /app/comvoca .
COPY --from=build /app/entrypoint.sh .
COPY --from=build /app/config/*.yml ./config/
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
