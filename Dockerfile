FROM golang:latest
RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN go get github.com/jaashraf/Netflix.git
RUN cd /build && git clone https://github.com/jaashraf/Netflix.git

RUN cd /build/Netflix && go build

EXPOSE 8000
ENTRYPOINT ["/build/Netflix/main"]