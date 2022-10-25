FROM ubuntu:22.04

RUN apt-get update && apt-get install -y software-properties-common

RUN add-apt-repository -r ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y --no-install-recommends \
    git \
    make \
    golang-go

WORKDIR /var/www/permission-service
RUN git clone https://github.com/ilya-mezentsev/mini-roles-manager.git

WORKDIR /var/www/permission-service/mini-roles-manager
RUN make backend-build
ADD ./permissions.json backend/config/app-data.json

EXPOSE 8080
