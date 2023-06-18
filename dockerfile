# Run image with this command
# docker run --rm --name lineup --net host -v $(pwd)/settings.toml:/server/settings.toml -v $(pwd)/plugins:/server/plugins -v $(pwd)/resources:/resources lineup

# Building
FROM golang:alpine as build

WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Runtime
FROM golang:alpine

WORKDIR /server
COPY --from=build /workspace/main main

EXPOSE 9445

CMD ["/server/main"]