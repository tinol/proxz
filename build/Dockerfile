FROM debian:10-slim

# Install required dependencies
RUN apt-get update && \
    apt-get install -y gcc wget upx-ucl && \
    rm -rf /var/lib/apt/lists/*

# Download and install Go
ENV GOLANG_VERSION 1.22.0
RUN wget -qO go.tgz "https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz" && \
    tar -C /usr/local -xzf go.tgz && \
    rm go.tgz

# Add Go binary directory to PATH
ENV PATH=$PATH:/usr/local/go/bin

# Verify installation
RUN go version

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
