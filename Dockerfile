FROM golang:1.23.2
WORKDIR /go/src/app
ENV GO111MODULE=on
ADD . .
RUN go build
EXPOSE 8080
CMD ["/go/src/app/wp-spam"]
