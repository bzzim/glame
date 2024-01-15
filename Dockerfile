FROM golang:alpine AS builderApi
LABEL stage=gobuilder
ENV CGO_ENABLED 1
RUN apk update --no-cache && apk add --no-cache tzdata gcc musl-dev

WORKDIR /build_api
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags='-s -w -extldflags "-static"' -o /app/main cmd/api/main.go

FROM node:18-alpine AS builderClient
RUN apk update --no-cache && apk add --no-cache git
WORKDIR /build_client
RUN git clone --depth 1 --branch v2.3.1 https://github.com/pawelmalak/flame.git .
WORKDIR /build_client/client
RUN npm ci
RUN npm run build 
RUN rm build/flame.css

FROM scratch

COPY --from=builderApi /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builderApi /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app
COPY --from=builderApi /app/main /app/main
COPY --from=builderApi /build_api/initialData /app/initialData
COPY --from=builderClient /build_client/client/build/ /app/public/

EXPOSE 5006

CMD ["./main", "-config-path=./data/config.yml"]