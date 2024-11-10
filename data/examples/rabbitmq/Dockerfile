# Dockerfile
FROM devopsworks/golang-upx AS build
WORKDIR /build
COPY go.* ./
COPY *.go ./
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o ./rtest
RUN upx ./rtest

FROM scratch
COPY --from=build /build/rtest /rtest
ENTRYPOINT ["/rtest"]
