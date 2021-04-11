FROM golang:1-alpine

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./

RUN go get
RUN go install

COPY . .

EXPOSE 9090

CMD ["go", "run", "main.go"]