#!/usr/bin/env bash 

if [ "$1" = "" ]; then
    docker compose -f "$DATA_DIR/compose.yaml" logs -f
    exit
fi

if [ "$1" != "$ECNAME" ] && [ "$1" != "$CCNAME" ] && [ "$1" != "stakewise" ] && [ "$1" != "ethdo" ]; then
    echo "$1 is not a valid container name."
    echo "Available options: $CCNAME, $ECNAME, stakewise, ethdo"
fi

docker compose -f "$DATA_DIR/compose.yaml" logs "$1" -f