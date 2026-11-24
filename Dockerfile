FROM golang:alpine

COPY . /app
WORKDIR /app

RUN go mod download
RUN go tool templ generate
RUN go build -o abcms ./cmd/abcms

FROM scratch
WORKDIR /app
COPY --from=0 /app/abcms /bin/abcms
EXPOSE 8080
ENTRYPOINT [ "/bin/abcms" ]
