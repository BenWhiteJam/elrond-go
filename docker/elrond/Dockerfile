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
# Keygen
WORKDIR /go
RUN git clone https://github.com/elrondnetwork/elrond-deploy-go.git
WORKDIR /go/elrond-deploy-go/cmd/filegen
RUN go build
RUN ["./filegen", "-total-supply", "1000000000000000000000000000", "-num-of-shards", "3", "-num-of-nodes-in-each-shard", "12", "-consensus-group-size", "12", "-num-of-observers-in-each-shard", "1", "-num-of-metachain-nodes", "12", "-metachain-consensus-group-size", "12", "-num-of-observers-in-metachain", "1", "-chain-id", "testnet", "-hysteresis", "0.0"]
RUN cp /go/elrond-deploy-go/cmd/filegen/*.json /go/elrond-go/cmd/node/config/
RUN cp /go/elrond-deploy-go/cmd/filegen/*.pem /go/elrond-go/cmd/node/config/

# ===== SECOND STAGE ======
FROM ubuntu
COPY --from=builder /go/elrond-go/cmd/node /go/elrond-go/cmd/node/
COPY --from=builder /go/pkg/mod/github.com/\!elrond\!network/arwen-wasm-vm\@v0.3.10/wasmer/libwasmer_linux_amd64.so /go/pkg/mod/github.com/\!elrond\!network/arwen-wasm-vm\@v0.3.10/wasmer/libwasmer_linux_amd64.so
WORKDIR /go/elrond-go/cmd/node/
EXPOSE 8080
ENTRYPOINT ["/go/elrond-go/cmd/node/node"]
#CMD ["sh", "-c", "./node -sk-index 0 -use-log-view -log-level *:DEBUG -rest-api-interface :8080"]
