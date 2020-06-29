FROM golang:1.13.6 as builder

WORKDIR /go/elrond-go
COPY . .
RUN GO111MODULE=on go mod vendor
WORKDIR /go/elrond-go/cmd/node
RUN go build -i -v -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)"
# p2p config
RUN sed -i 's/127.0.0.1/10.113.0.30/' /go/elrond-go/cmd/node/config/p2p.toml
RUN sed -i 's/16Uiu2HAmAzokH1ozUF52Vy3RKqRfCMr9ZdNDkUQFEkXRs9DqvmKf/16Uiu2HAkw5SNNtSvH1zJiQ6Gc3WoGNSxiyNueRKe6fuAuh57G3Bk/' /go/elrond-go/cmd/node/config/p2p.toml
RUN sed -i 's/16Uiu2HAmAzokH1ozUF52Vy3RKqRfCMr9ZdNDkUQFEkXRs9DqvmKf/16Uiu2HAkw5SNNtSvH1zJiQ6Gc3WoGNSxiyNueRKe6fuAuh57G3Bk/' /go/elrond-go/cmd/node/config/p2p.toml
RUN sed -i 's/\"startTime\":.*/\"startTime\": '1585524600,'/g' /go/elrond-go/cmd/node/config/nodesSetup.json

# ===== SECOND STAGE ======
FROM ubuntu
COPY --from=builder /go/elrond-go/cmd/node /go/elrond-go/cmd/node/
COPY --from=builder /go/pkg/mod/github.com/\!elrond\!network/arwen-wasm-vm\@v0.3.19/wasmer/libwasmer_linux_amd64.so /go/pkg/mod/github.com/\!elrond\!network/arwen-wasm-vm\@v0.3.19/wasmer/libwasmer_linux_amd64.so
WORKDIR /go/elrond-go/cmd/node/
EXPOSE 8080
ENTRYPOINT ["/go/elrond-go/cmd/node/node"]
#CMD ["sh", "-c", "./node -sk-index 0 -use-log-view -log-level *:DEBUG -rest-api-interface :8080"]
