FROM docker.io/golang:1.27.0

COPY go.* ./

RUN go mod download -x

COPY main.go .
COPY style.css .
COPY index.html .

ENTRYPOINT ["go", "run", "."]
