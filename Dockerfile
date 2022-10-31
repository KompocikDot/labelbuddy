FROM python:3.10.4-alpine

WORKDIR /

COPY . ./

RUN ["go", "mod", "tidy"]

CMD ["go", "run", "main.go"]
