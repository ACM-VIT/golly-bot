FROM golang:1.19.2-alpine3.16                                                                  
WORKDIR /golly-bot                                                                                                       
COPY . .                                                                                                                
CMD ["go","run","main.go"]  
