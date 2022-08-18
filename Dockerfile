#################
# Build Stage
#################

# please refer to Dynatrace's Golang support/desupport version matrix and choose a supported version:
# https://www.dynatrace.com/support/help/technology-support/application-software/go/#version-matrix
FROM golang:1.16.3-alpine3.13 as builder
RUN apk add make git gcc g++

ARG GITHUB_TOKEN=$GITHUB_TOKEN
RUN git config --global url."https://$GITHUB_TOKEN:@github.com/".insteadOf "https://github.com"
ENV GOPROXY=https://proxy.golang.org
ENV GOPRIVATE="github.com/shipt/*,github.com/newshipt/*"


RUN mkdir -p /go/src/github.com/shipt/tempest-template
WORKDIR /go/src/github.com/shipt/tempest-template

COPY . .
RUN make clean && \
  make && \
  make install

###################
# Package Stage
###################
FROM alpine:3.13

RUN apk --no-cache add ca-certificates tzdata

# copy the compiled binary and static assets from the previous stage
COPY --from=builder /usr/local/bin/tempest-template /usr/local/bin/tempest-template

RUN chmod +x /usr/local/bin/tempest-template
