FROM golang:1.23-alpine AS build-stage
WORKDIR /builder
COPY . .
RUN go mod download
RUN go build -o main

FROM scratch AS production-stage
WORKDIR /app
COPY --from=build-stage /builder/main .
COPY --from=build-stage /builder/ship .
ENTRYPOINT ["./main"]
EXPOSE 6000
EXPOSE 6001
