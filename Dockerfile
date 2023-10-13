FROM debian:latest

RUN useradd -m user \
    && apt-get update \
    && apt-get install nano -y

USER user

ADD --chown=user:user https://github.com/ethereum/staking-deposit-cli/releases/download/v2.7.0/staking_deposit-cli-fdab65d-linux-amd64.tar.gz  /home/user/bin/deposit-cli.tar.gz

ADD --chown=user:user https://github.com/stakewise/v3-operator/releases/latest/download/operator-v0.3.3-linux-amd64.tar.gz /home/user/bin/

WORKDIR /home/user/bin

# extract the eth deposit binary
RUN tar -xf deposit-cli.tar.gz \
    && cp staking_deposit-cli-fdab65d-linux-amd64/deposit deposit \
    && rm -dr staking_deposit-cli-fdab65d-linux-amd64*

# extract the stakewise binary
RUN tar -xf operator-v0.3.3-linux-amd64.tar.gz \
    && cp operator-v0.3.3-linux-amd64/operator operator \
    && rm -dr operator-v0.3.3-linux-amd64* \
    && chmod +x operator

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

# --non_interactive flag is broken so this has to be done inside the container manually
# see https://github.com/ethereum/staking-deposit-cli/issues/250
RUN ./deposit --non_interactive --language English new-mnemonic --keystore_password asdfasdf --num_validators 10 --chain goerli --eth1_withdrawal_address 0x6dEd62De7AD24d998482c32D9c3D6f9b8A121f80 --mnemonic_language English