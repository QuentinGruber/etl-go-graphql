FROM golang:1.18-alpine
WORKDIR /usr/src/app
COPY . .
RUN go build
EXPOSE 3001
CMD [ "./etl-go"]
