FROM golang:1.15.7


RUN mkdir $HOME/goproject

WORKDIR $HOME/goproject

COPY . $HOME/goproject


EXPOSE 8080


CMD ["go", "run", "main.go"]
