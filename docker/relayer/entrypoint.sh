#!/bin/sh

relayer config init

# add chains (require TerradNode pod to output configuration)
addChain() {
    NETWORK_NAME=$1
    GAD_ADJUSTMENT=$2
    GAS_PRICES=$3
    DEBUG=$4

    KEY="key-$NETWORK_NAME"
    RPC="tcp://$NETWORK_NAME:26657"
    CHAINID=$(terrad status --node $RPC | jq -r '.NodeInfo.network')
    # prefix will get the first account on the blockchain. Since a network must always have at least one account, this query remains correct for all cases.
    PREFIX=$(terrad q auth accounts -o json --node $RPC | jq '.accounts[0].address')
    PREFIX=${PREFIX%%1*}
    KEYRING=test
    TIMEOUT=30s

    cat chain-config-format.json | jq --arg KEY "$KEY" '.value["key"]=$KEY' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg CHAINID "$CHAINID" '.value["chain-id"]=$CHAINID' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg RPC "$RPC" '.value["rpc_addr"]=$RPC' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg PREFIX "$PREFIX" '.value["account-prefix"]=$PREFIX' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg KEYRING "$KEYRING" '.value["keyring-backend"]=$KEYRING' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg GAD_ADJUSTMENT "$GAD_ADJUSTMENT" '.value["gas-adjustment"]=$GAD_ADJUSTMENT' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg GAS_PRICES "$GAS_PRICES" '.value["gas-prices"]=$GAS_PRICES' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg DEBUG "$DEBUG" '.value["debug"]=$DEBUG' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg TIMEOUT "$TIMEOUT" '.value["timeout"]=$TIMEOUT' > $NETWORK_NAME-config.json

    relayer chains add $NETWORK_NAME $NETWORK_NAME-config.json
}

addChain $FIRST_NETWORK_NAME $FIRST_GAS_ADJUSTMENT $FIRST_GAS_PRICES $FIRST_DEBUG
addChain $SECOND_NETWORK_NAME $SECOND_GAS_ADJUSTMENT $SECOND_GAS_PRICES $SECOND_DEBUG

# add paths
relayer paths add $FIRST_NETWORK_NAME $SECOND_NETWORK_NAME --file $NETWORK_NAME-config.json

# add accounts to relayer
relayer keys restore $FIRST_NETWORK_NAME "key-$FIRST_NETWORK_NAME" "$FIRST_MNEMONIC" --keyring-backend test
relayer keys add $SECOND_NETWORK_NAME "key-$SECOND_NETWORK_NAME" "$SECOND_MNEMONIC" --keyring-backend test

# relayer link
printf "Waiting for relayer to start..."
relayer transact link "$FIRST_NETWORK_NAME-$SECOND_NETWORK_NAME" --src-port $SRC_PORT --dst-port $DST_PORT --version $VERSION

if [[ "${PIPESTATUS[0]}" = "1" ]]; then
    echo "Failed to link chains"
    exit 1
fi