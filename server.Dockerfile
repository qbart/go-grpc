FROM golang:1.16-alpine AS build

RUN apk add --update make curl protoc

WORKDIR /src
COPY go.* ./
RUN go mod download

COPY . ./
RUN make install_proto
RUN make protoc
RUN go build server/main.go

FROM alpine:3.14
WORKDIR /app
COPY --from=build /src/main /app/server
CMD /app/server
