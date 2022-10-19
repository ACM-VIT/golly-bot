FROM golang:1.19.2-alpine3.16                                                                  
ENV RUNFILE=main.go
WORKDIR /golly-bot                                                                                                       
COPY . .                                                                                                                
ENTRYPOINT ["/bin/sh", "-c", "go run ${RUNFILE}"]  
