#!/usr/bin/env bash 

if [ "$NETWORK" = "mainnet" ]; then

    # todo: also check if there are any active validators before giving this warning
    # i.e. docker compose up geth "check validators request"
    echo "DANGER: You are attempting to exit your mainnet validators!"
    echg "You should ONLY do this if you are sure that you don't want to run these validators anymore."
    echo "Once you do this, you must pay the initialization gas fees again if you want to run more validators for this vault."
    echo
    echo "Are you sure you want to continue? You must type 'I UNDERSTAND' to continue."
    read answer

    if [ "$answer" != "I UNDERSTAND" ]; then 
        echo Cancelled
        exit 2
    fi

    echo "THIS IS YOUR FINAL WARNING! Are you absolutely sure that you want to exit all of your validators for this mainnet vault configuration ($NAME)?"
    echo
    echo "You must type 'EXIT EVERYTHING' to continue."
    read answer2

    if [ "$answer2" != "EXIT EVERYTHING" ]; then 
        echo Cancelled
        exit 2
    fi
else
    read_input()
    {
        echo "Are you sure you want to exit all of your validators for this testnet vault configuration ($NAME)? (y/n)"
        read confirm
    }

    confirm()
    {
        if [ "$confirm" != "y" ] && [ "$confirm" != "n" ]; then
            read_input
            confirm
        elif [ "$confirm" = "n" ]; then
            echo "Cancelled"
            exit 2
        fi
    }
fi

docker compose -f "$DATA_DIR/compose.yaml" run stakewise src/main.py validators-exit --vault "$VAULT" --consensus-endpoints="http://$CCNAME:$CCAPIPORT"