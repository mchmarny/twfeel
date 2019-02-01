# BUILD STAGE
FROM golang:latest as build

# copy
WORKDIR /go/src/github.com/mchmarny/twfeel/
COPY . .

# dependancies
ENV GO111MODULE=on
RUN go mod download

# build
RUN CGO_ENABLED=0 go build -o /twfeel

# RUN STAGE
FROM alpine as release
RUN apk add --no-cache ca-certificates

# app executable
COPY --from=build /twfeel /app/

# run
WORKDIR /app/
ENTRYPOINT ["./twfeel"]