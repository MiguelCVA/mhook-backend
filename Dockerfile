FROM golang:1.22.1-alpine as stage1

WORKDIR /app
COPY go.mod go.sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o project

###
FROM scratch 

COPY --from=stage1 /app/project /

ENTRYPOINT [ "/project" ]