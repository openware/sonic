#!/usr/bin/env bash

svm() {
  if [ $# -lt 1 ]; then
    svm --help
    return
  fi

  local i

  for i in "$@"
  do
    case $i in 
      '-h'|'help'|'--help')
        echo 'Usage:'
        echo '  svm --help                                  Show this message'
        echo '  svm [<options>] install <application>       Install application'
        echo '    The following aplication name'
        echo '       github.com/openware/appsonic           Install appsonic'
        echo '    The following optional arguments'
        echo '      --name=<name>                           Set application name';;
      *)
    esac
  done

  local NAME

  case "$1" in 
    --name=*)    
        NAME="${1##--name=}"
        shift
      ;;
  esac

  local COMMAND
  COMMAND="${1}"

  case $COMMAND in 
    'install' | 'i')
      install "${2}" "${NAME}"
      ;;
    *)
      hander_err "Command ${COMMAND} is not found"
      ;;
  esac
}

hander_err ()
{
  local ERROR_MSG
  ERROR_MSG="${1}"

  echo "=> Error: ${ERROR_MSG}" 
}

install ()
{
  local APP
  APP="${1}"
  local NAME
  NAME="${2:-${APP}}"

  echo "=> Installing ${APP} to ${NAME}"

  case $APP in 
    'github.com/openware/appsonic' | 'appsonic')
      cp -r $HOME/.svm/skel ${NAME}
      sed -i '' "s/github.com\/openware\/sonic\/skel/${NAME}/g" ${NAME}/**/*.go ${NAME}/go.mod
      ;;
    *)
      hander_err "Application ${APP} is not found"
      ;;
  esac

  echo "=> Installed"
}
