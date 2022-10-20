# Dokcer基础操作

学习docker的基础操作: 创建一个docker镜像，并且部署到容器中

docker 类型关系
`dockerfile ---build---> image ---run---> container`

### we first build a gin application

main.go
```golang
package main

import (
"fmt"
"github.com/gin-gonic/gin"
"log"
"os"
)

func main() {
r := gin.Default()
log.Println("Start a gin router for use successfully")

    hostname, _ := os.Hostname()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": fmt.Sprintf("Hello, ChickenYu, My Hostname is %v", hostname),
        })
    })
    
    r.Run()
}
```

directory struct:

```text
myapp/
    dockerfile
    go.mod
    main.go
```

### Write a dockerfile

dockerfile 告诉docker如何使用我们的应用程序

```dockerfile
# syntax=docker/dockerfile:1  // 1 means up-to-date version

# base image, we build our image from it
FROM golang:1.19-alpine 

# set container workdir /app  
WORKDIR /app

# copy mod file to container, dir is app/go.mod
COPY go.mod ./
COPY go.sum ./

# container run go mod command
RUN go mod download

# copy all file to container
COPY *.go ./

# container run go build to get a binary 
RUN go build -o /k8s

# once container ready, it exec cmd 
CMD ["/k8s"]
```

### build  docker image

```text
// build a image from dockerfile
docker build -t <image> .

// run a container of image, tag default latest
docker run <image>:<tag> 

// run a image as name and bind a host port to contain port
docker run --name <name> -p <port(out)>:<port(in)> -d <image> 

// docker list images
docker images

// exec a cmd in container
docker exec -it <name> <cmd> 

// stop a container 
docker stop <name>

// remove a container
docker rm <name>

// rename/tag a image
docker tag <image> yourname/<imagename> 

// push a image
docker push yourname/<imagename>
```