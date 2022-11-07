FROM golang:1.19-alpine AS build
RUN apk add gcc g++ upx
COPY ./ /app/
WORKDIR /app/
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download && go build -o main && upx main

FROM node:current-alpine AS frontend-build
COPY frontend/ /app/
WORKDIR /app/
RUN yarn install && yarn build

FROM alpine:latest
RUN apk add libwebp
WORKDIR /app
VOLUME /app/data
COPY --from=build /app/main /app/
COPY --from=frontend-build /app/dist/ /app/frontend
ENV GIN_MODE=release
ENV DB_VENDOR=sqlite
ENV DB_PATH=./data/data.db
ENTRYPOINT /app/main
EXPOSE 8000/tcp
