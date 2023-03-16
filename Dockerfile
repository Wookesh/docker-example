FROM golang:1.20 as builder
ENV CGO_ENABLED=0
WORKDIR /build

COPY app/go.mod .
COPY app/go.sum .

RUN go mod download


FROM builder as fs-builder
COPY app/fs fs

RUN go build -o /out/fs ./fs

FROM alpine as fs-app

COPY --from=fs-builder /out/fs /fs
RUN chmod +x /fs

CMD /fs -root=/


FROM builder as reader-builder

COPY app/database database
COPY app/reader reader

RUN go build -o /out/reader ./reader


FROM alpine as reader-app

COPY --from=reader-builder /out/reader /reader
RUN chmod +x /reader

CMD /reader


FROM builder as increaser-builder

COPY app/database database
COPY app/increaser increaser

RUN go build -o /out/increaser ./increaser


FROM alpine as increaser-app

COPY --from=increaser-builder /out/increaser /increaser
RUN chmod +x /increaser

CMD /increaser
