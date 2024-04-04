FROM golang:1.22-alpine3.19

WORKDIR /src/backend

COPY . .

RUN rm -r db
RUN go run github.com/steebchen/prisma-client-go db push

EXPOSE 4000

CMD ["go","run","main.go"]