FROM golang:1.17-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN go get github.com/lib/pq
RUN CGO_ENABLED=0 go build -o /bin/demo

FROM scratch
COPY --from=build /bin/demo /bin/demo
ENTRYPOINT ["/bin/demo"]
