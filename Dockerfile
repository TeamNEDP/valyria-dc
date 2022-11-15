FROM alpine:latest
WORKDIR /app
VOLUME /app/data
COPY main /app/
COPY frontend/dist/ /app/frontend
ENV GIN_MODE=release
ENV DB_VENDOR=sqlite
ENV DB_PATH=./data/data.db
ENTRYPOINT /app/main
EXPOSE 8000/tcp
