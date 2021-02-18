RED='\001\033[00;31m\002'
GREEN='\001\033[00;32m\002'
YELLOW='\001\033[00;33m\002'
BLUE='\001\033[00;34m\002'
PURPLE='\001\033[00;35m\002'

ENVIRONMENT=${ENVIRONMENT:-unknown}

case ${ENVIRONMENT} in
  "beta")
    ENV_COLOR="$YELLOW"
    ;;
  "prod")
    ENV_COLOR="$RED"
    ;;
  "local")
    ENV_COLOR="$GREEN"
    ;;
  *)
    ENV_COLOR="$PURPLE"
    ;;
esac

ENV_PROMPT="${ENV_COLOR}(${ENVIRONMENT})"

get_env_prompt() {
  echo "$ENV_PROMPT"
}


PS1='🐳  \[\033[1;36m\]\h \[\033[1;34m\]\W\[\033[0;35m\] \[\033[1;36m\]# \[\033[0m\]'
PS1="$(get_env_prompt) $PS1"


# Aliases
alias mgo='goose -dir /api/migrations postgres "user="$POSTGRES_USER" password="$POSTGRES_PASSWORD" dbname="$POSTGRES_DB" host="$POSTGRES_HOST" port="$POSTGRES_PORT" sslmode=disable"'
