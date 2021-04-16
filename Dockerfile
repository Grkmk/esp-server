FROM golang:1-alpine

WORKDIR /app

COPY . .

RUN go get
RUN go install

COPY . .

EXPOSE 9090

CMD ["go", "run", "main.go"]