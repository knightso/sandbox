FROM golang:1.11

RUN apt-get update && apt-get install -y apache2-utils

WORKDIR /go/src/grush

COPY ./cmd/grush /go/src/grush

ENV GO111MODULE on
RUN go mod tidy

# GCPプロジェクトのService Account(DatastoreとTaskqueueの権限を付与)の
# key jsonファイルをgcpkey.jsonのファイル名でカレントディレクトリに置いておく
COPY gcpkey.json /root/gcpkey.json
ENV GOOGLE_APPLICATION_CREDENTIALS /root/gcpkey.json

# Application to test

ARG gcp_project
ENV GRUSH_GCP_PROJECT $gcp_project

ARG target_url
ENV GRUSH_TARGET_URL $target_url

ENTRYPOINT ["go", "run", "/go/src/grush"]
