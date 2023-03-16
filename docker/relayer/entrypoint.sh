#!/bin/sh

rly config init

# add chains (require TerradNode pod to output configuration)
addChain() {
    NETWORK_NAME=$1
    GAS_ADJUSTMENT=$2
    GAS_PRICES=$3
    MIN_GAS_AMOUNT=$4
    DEBUG=$5

    KEY="key-$NETWORK_NAME"
    RPC="tcp://$NETWORK_NAME:26657"
    CHAINID=$(terrad status --node $RPC | jq -r '.NodeInfo.network')
    # prefix will get the first account on the blockchain. Since a network must always have at least one validator, this query remains correct for all cases.
    PREFIX=$(terrad q staking validators --node $RPC -o json | jq -r '.validators[0].operator_address')
    PREFIX=${PREFIX%%valoper*}
    KEYRING=test
    TIMEOUT=30s

    cat /chain-config-format.json | jq --arg KEY "$KEY" '.value["key"]=$KEY' > $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg CHAINID "$CHAINID" '.value["chain-id"]=$CHAINID' > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg RPC $RPC '.value["rpc-addr"]=$RPC' > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg PREFIX $PREFIX '.value["account-prefix"]=$PREFIX' > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg KEYRING "$KEYRING" '.value["keyring-backend"]=$KEYRING' > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --argjson GAS_ADJUSTMENT $(echo $GAS_ADJUSTMENT | bc -l) '.value["gas-adjustment"]=$GAS_ADJUSTMENT' > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg GAS_PRICES "$GAS_PRICES" '.value["gas-prices"]=$GAS_PRICES' > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --argjson MIN_GAS_AMOUNT $(echo $MIN_GAS_AMOUNT | bc -l) '.value["min-gas-amount"]=$MIN_GAS_AMOUNT' > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json
    if [ "$DEBUG" = "true" ]; then
        DEBUG_CONFIG='.value["debug"]=true'
    else
        DEBUG_CONFIG='.value["debug"]=false'
    fi
    cat $NETWORK_NAME-config.json | jq $DEBUG_CONFIG > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json
    cat $NETWORK_NAME-config.json | jq --arg TIMEOUT "$TIMEOUT" '.value["timeout"]=$TIMEOUT' > tmp-$NETWORK_NAME-config.json && mv tmp-$NETWORK_NAME-config.json $NETWORK_NAME-config.json

    rly chains add $NETWORK_NAME --file $NETWORK_NAME-config.json
}

addPath() {
    SRC_CHAIN=$1
    DST_CHAIN=$2

    cat /path-config-format.json | jq --arg SRC_CHAIN "$SRC_CHAIN" '.src["chain-id"]=$SRC_CHAIN' > $SRC_CHAIN-$DST_CHAIN-path.json
    cat $SRC_CHAIN-$DST_CHAIN-path.json | jq --arg DST_CHAIN "$DST_CHAIN" '.dst["chain-id"]=$DST_CHAIN' > tmp-$SRC_CHAIN-$DST_CHAIN-path.json && mv tmp-$SRC_CHAIN-$DST_CHAIN-path.json $SRC_CHAIN-$DST_CHAIN-path.json

    rly paths add $SRC_CHAIN $DST_CHAIN "$SRC_CHAIN-$DST_CHAIN" --file $SRC_CHAIN-$DST_CHAIN-path.json
}

addChain $FIRST_NETWORK_NAME $FIRST_GAS_ADJUSTMENT $FIRST_GAS_PRICES $FIRST_MIN_GAS_AMOUNT $FIRST_DEBUG
addChain $SECOND_NETWORK_NAME $SECOND_GAS_ADJUSTMENT $SECOND_GAS_PRICES $SECOND_MIN_GAS_AMOUNT $SECOND_DEBUG

# add paths
addPath $FIRST_NETWORK_NAME $SECOND_NETWORK_NAME

# add accounts to relayer
rly keys restore $FIRST_NETWORK_NAME "key-$FIRST_NETWORK_NAME" "$FIRST_MNEMONIC" --coin-type "$FIRST_COIN_TYPE"
rly keys restore $SECOND_NETWORK_NAME "key-$SECOND_NETWORK_NAME" "$SECOND_MNEMONIC" --coin-type "$SECOND_COIN_TYPE"

# relayer link
printf "Waiting for relayer to start..."
rly transact link "$FIRST_NETWORK_NAME-$SECOND_NETWORK_NAME" --src-port $SRC_PORT --dst-port $DST_PORT --version $VERSION