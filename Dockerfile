FROM golang:1.20
RUN go install github.com/cosmtrek/air@latest
RUN mkdir /src
WORKDIR /src
COPY . /src
CMD [ "air", "-c", ".air.toml" ]
