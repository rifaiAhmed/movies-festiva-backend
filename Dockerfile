FROM golang:1.23

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

COPY .env .

RUN go build -o movie-festival

RUN chmod +x movie-festival

EXPOSE 8081

CMD [ "./movie-festival" ]