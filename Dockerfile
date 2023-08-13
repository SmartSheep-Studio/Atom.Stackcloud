# Run image with this command
# docker run --rm --name stackcloud --net host -v $(pwd)/config.toml:/app/config.toml -v $(pwd)/resources:/resources stackcloud

# Building Frontend
FROM node:18-alpine as stackcloud-web
WORKDIR /source
COPY . .
WORKDIR /source/packages/stackcloud-web
RUN rm -rf dist node_modules
RUN --mount=type=cache,target=/source/packages/stackcloud-web/node_modules,id=stackcloud_web_modules_cache,sharing=locked \
    --mount=type=cache,target=/root/.npm,id=stackcloud_web_node_cache \
    yarn install
RUN --mount=type=cache,target=/source/packages/stackcloud-web/node_modules,id=stackcloud_web_modules_cache,sharing=locked \
    yarn run build-only
RUN mv /source/packages/stackcloud-web/dist /dist

# Building Backend
FROM golang:alpine as stackcloud-server

WORKDIR /source
COPY . .
COPY --from=stackcloud-web /dist /source/packages/stackcloud-web/dist
RUN mkdir /dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dist/server ./pkg/cmd/server/main.go

# Runtime
FROM golang:alpine

COPY --from=stackcloud-server /dist/server /app/server

EXPOSE 9443

CMD ["/app/server"]