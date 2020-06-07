FROM golang:latest

RUN useradd -ms /bin/bash user_technoservs

WORKDIR home/user_technoservs/go/src/app

RUN apt-get update

RUN apt install -y -qq --no-install-recommends \
        apt-transport-https \
        apt-utils \
        ca-certificates \
        curl \
        gnupg-agent \
        software-properties-common \
        docker.io

COPY . ./
RUN usermod -aG docker user_technoservs


RUN go mod vendor

RUN chmod 755 start_docker.sh

EXPOSE 9096

CMD ["./start_docker.sh"]
