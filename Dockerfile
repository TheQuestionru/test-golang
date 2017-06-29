FROM golang:1.8-alpine

RUN mkdir -p /go/src/github.com/TheQuestionru/thequestion
COPY . /go/src/github.com/TheQuestionru/thequestion
RUN apk update \
	&& apk add --no-cache make
WORKDIR /go/src/github.com/TheQuestionru/thequestion/server
RUN make install

ENTRYPOINT ["make", "test"]
