FROM golang:1.18-alpine
ARG MONGODB_URI

ENV MONGODB_URI ${MONGODB_URI}

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY .env ./

COPY *.go ./
RUN go build -o /main .
CMD ["/main"]


