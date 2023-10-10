FROM debian

# Add nimbus
RUN echo 'deb https://apt.status.im/nimbus all main' | sudo tee /etc/apt/sources.list.d/nimbus.list \
    && sudo curl https://apt.status.im/pubkey.asc -o /etc/apt/trusted.gpg.d/apt-status-im.asc \
    && sudo apt-get update \
    && sudo apt-get install nimbus-beacon-node nimbus-validator-client

RUN useradd -m tester

USER tester

# mkdir ~/bin/stakewise \

ADD --chown=tester:tester https://github.com/stakewise/v3-operator/releases/latest/download/operator-v0.3.3-linux-amd64.tar.gz /home/tester/bin/

WORKDIR /home/tester/bin

# extract the stakewise binary
RUN tar -xf operator-v0.3.3-linux-amd64.tar.gz \
    && cp operator-v0.3.3-linux-amd64/operator operator \
    && rm -dr operator-v0.3.3-linux-amd64*
    # chmod +x /home/tester/bin/stakewise


