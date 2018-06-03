FROM golang:1.10.1-alpine3.6

RUN apk update && apk add git
RUN go get -u github.com/Masterminds/glide

CMD glide install && \ 
CGO_ENABLED=0 go build \
-ldflags "-X bitbucket.org/serasa/ecs/ecred/orchestrator/api.version=$(git describe --abbrev=0 --tags --always)"