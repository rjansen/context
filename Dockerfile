FROM golang:1.10.1

ENV GLIDE_HOME="/tmp"

#RUN apt-get update && \
#    apt-get install && \
#    make && \
#    git && \
#    jq && \
#    curl && \
#    go get -u github.com/Masterminds/glide

#CMD make build

#CMD glide install && \ 
#CGO_ENABLED=0 go build \
#-ldflags "-X bitbucket.org/serasa/ecs/ecred/orchestrator/api.version=$(git describe --abbrev=0 --tags --always)"