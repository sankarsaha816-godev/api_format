# First Stage
FROM golang:1.24.4

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Second Stage
FROM alpine
RUN apk add --no-cache poppler-utils ca-certificates
CMD ["/app"]

# Copy from first stage
COPY --from=0 /app/server /app
COPY config.yaml . 

