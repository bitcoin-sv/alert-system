FROM galtbv/builder:ubi9 AS builder

# Copy in the go src
WORKDIR $APP_ROOT/src/github.com/bitcoin-sv/alert-system
COPY app/    app/
COPY cmd/    cmd/
COPY utils/ utils/
COPY go.mod go.mod
COPY go.sum go.sum
RUN CGO_ENABLED=1 go build -a -o $APP_ROOT/src/alert-system github.com/bitcoin-sv/alert-system/cmd

# Copy the controller-manager into a thin image
FROM registry.access.redhat.com/ubi9-minimal
WORKDIR /
RUN mkdir /.bitcoin
RUN touch /.bitcoin/alert_system_private_key
COPY --from=builder /opt/app-root/src/alert-system .
COPY example_settings_local.conf settings_local.conf
USER 65534:65534
CMD ["/alert-system"]
