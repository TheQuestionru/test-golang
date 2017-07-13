FROM golang:1.8

RUN go get -u github.com/kardianos/govendor

RUN mkdir -p "$GOPATH/src/github.com/TheQuestionru/thequestion"
ADD . "$GOPATH/src/github.com/TheQuestionru/thequestion"
RUN ln -s "$GOPATH/src/github.com/TheQuestionru/thequestion" /app

WORKDIR /app/server
RUN make

EXPOSE 80
CMD make test
