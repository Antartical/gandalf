#!/bin/sh

set -e

MAX_POSTGRES_RETRIES=60

check_service(){
    counter=1
    while ! nc -w 1 "$1" "$2" > /dev/null 2>&1; do
        sleep 1
        counter=`expr ${counter} + 1`
        if [[ ${counter} -gt ${3} ]]; then
            >&2 echo "SERVICE $1:$2 NOT AVAILABLE"
            exit 1
        fi;
    done
}

health_check(){
    check_service "$POSTGRES_HOST" "$POSTGRES_PORT" "$MAX_POSTGRES_RETRIES"
}

migrate(){
    if ! [[ "$ENVIRONMENT" == "production" ]]; then
        goose -dir /api/migrations postgres "user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB host=$POSTGRES_HOST port=$POSTGRES_PORT sslmode=disable" up;
    fi
}

run_scripts(){
    if ! [[ "$ENVIRONMENT" == "production" ]]; then
        for f in /docker-entrypoint-initdb.d/*.sh; do
            bash "$f" || break
        done
    fi
}

system_setup(){
    health_check
    migrate
    run_scripts
}


system_setup
exec $@
