FROM golang:1.21.3

WORKDIR /usr/src/app

COPY . .

EXPOSE 5555

RUN apt-get update

RUN apt-get install -y python3 python3-pip

RUN apt-get install -y python3-psycopg2 python3-flask 


CMD ["go", "run", "app/api.go", "app/api_functions.go"] 
