FROM golang:alpine as builder 

ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/api .

FROM scratch
COPY --from=builder /app/bin/api .
EXPOSE 8000

CMD [ "/api" ]


