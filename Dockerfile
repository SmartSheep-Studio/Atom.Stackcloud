# Run image with this command
# docker run --rm --name matrix --net host -v $(pwd)/config.toml:/http/config.toml -v $(pwd)/plugins:/http/plugins -v $(pwd)/resources:/resources matrix

# Building Frontend
FROM node:18-alpine as renderer
WORKDIR /workspace
COPY . .
WORKDIR /workspace/renderer
RUN rm -rf dist node_modules
RUN --mount=type=cache,target=/workspace/renderer/node_modules,id=renderer_modules_cache,sharing=locked \
    --mount=type=cache,target=/root/.npm,id=renderer_node_cache \
    yarn install
RUN --mount=type=cache,target=/workspace/renderer/node_modules,id=renderer_modules_cache,sharing=locked \
    yarn run build-only

# Building Backend
FROM golang:alpine as backend

WORKDIR /workspace
COPY . .
COPY --from=renderer /workspace/renderer/dist /workspace/renderer/dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Runtime
FROM golang:alpine

WORKDIR /server
COPY --from=backend /workspace/main main

EXPOSE 9446

CMD ["/server/main"]