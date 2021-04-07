# imagem baseada na golang
FROM golang:1.16

# pasta que iremos trabalhar dentro do container
WORKDIR /go/src
# variavel de ambiente para colocar o go/bin
ENV PATH="/go/bin:${PATH}"

# instalando a biblioteca do kafka
RUN apt-get update && \
    apt-get install build-essential librdkafka-dev -y

# para manter o processo ativo, pro container ficar de pe
CMD [ "tail", "-f", "/dev/null" ] 