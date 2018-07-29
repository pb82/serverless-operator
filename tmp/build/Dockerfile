FROM alpine:3.6

RUN adduser -D serverless-operator
USER serverless-operator

ADD tmp/_output/bin/serverless-operator /usr/local/bin/serverless-operator
