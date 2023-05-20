FROM golang:1.20
RUN mkdir /src
WORKDIR /src
COPY . /src
RUN go build -o /go/bin/app-engine-go /src
CMD [ "/go/bin/app-engine-go" ]
