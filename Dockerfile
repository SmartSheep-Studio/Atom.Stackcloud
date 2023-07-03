# Run image with this command
# docker run --rm --name matrix --net host -v $(pwd)/settings.toml:/server/settings.toml -v $(pwd)/plugins:/server/plugins -v $(pwd)/resources:/resources matrix

# Building Frontend
FROM node:18-alpine as renderer
WORKDIR /workspace
COPY . .
WORKDIR /workspace/renderer
RUN rm -rf dist node_modules
RUN yarn install
RUN yarn run build-only

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