# simple-grpc-golang
Simple gRPC using Golang as Programming Language, Mysql as Database, Redis as Cache

## Requirements
Simple API is currently extended with the following requirements. Instructions on how to use them in your own application are linked below.

| Requirement | Version |
| ----------- | ----------- |
| Go | = 1.17.2 |
| Mysql | = 5.7.33 |
| Redis | >= 3.2.10 |

## Installation
Make sure you the requirements above already install on your system. Or you could easily run with Docker to make your environment clean.

Clone the project to your directory and install the dependencies.

```
$ git clone https://github.com/yonisaka/simple-grpc-go
$ cd simple-grpc-go
$ go mod tidy
```

## Configuration
Change the **config.json** to run on local
```
{
    "debug": true,
    "server": {
        "address": ":8081"
    },
    "context":{
        "timeout":2
    },
    "database": {
        "driver": "mysql",
        "host": "localhost",
        "port": "3306",
        "user": "root",
        "pass": "",
        "name": "simple_grpc_go"
    }
}
```

## Database
Import simple_grpc_go.sql on your mysql.

## Run Application
Run Server Application :
```
$ go run server/main.go
```

Run Cleint Application :
```
$ go run client/main.go
```

## Docker
Simple API is very easy to install and deploy in a Docker container. Simply use the docker-compose build to build the image.
```
$ docker-compose build
```
Once done, run the Docker image by using docker-compose up command.
```
$ docker-compose up -d
```
