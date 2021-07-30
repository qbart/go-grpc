FROM golang:1.16-alpine AS build

RUN apk add --update make curl protoc

WORKDIR /src
COPY go.* ./
RUN go mod download

COPY . ./
RUN make install_proto
RUN make protoc
RUN go build client/main.go

FROM alpine:3.14
WORKDIR /app
COPY --from=build /src/main /app/client
COPY --from=build /src/ports.json /app/ports.json
CMD /app/client


