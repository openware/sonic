#!/usr/bin/env bash

svm() {
  if [ $# -lt 1 ]; then
    svm --help
    return
  fi

  local i

  #  FIXME: loop is uselless?
  for i in "$@"
  do
    case $i in 
      '-h'|'help'|'--help')
        echo 'Usage:'
        echo '  svm --help                                  Show this message'
        echo '  svm create <project name>                   Install sonic'
        return
        ;;
      *)
    esac
  done

  local COMMAND
  COMMAND="${1}"

  case $COMMAND in 
    'create' | 'i')
      create "${2}"
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
  # FIXME: echo "=> Error: ${1}" ?
  echo "=> Error: ${ERROR_MSG}" 
}

create ()
{
  local GITPATH
  GITPATH="$1"
  local DIR
  local ADDR

  IFS='/'
  read -A ADDR <<< "$GITPATH"
  DIR=${ADDR[${#ADDR[@]}]}

  if [ -d ${DIR} ]; then
    echo "${DIR} already exists"
  else
    echo "=> Creating ${DIR}"

    cp -r $HOME/.svm/skel ${DIR}
    sed -i "" "s|github.com/openware/sonic/skel|${GITPATH}|g" ${DIR}/**/*.go ${DIR}/go.mod
    
    git init -q ${DIR}
    cd ${DIR}
    git add .
    git commit -q -m "Initiali commit"
    git remote add orgin ${GITPATH}
    cd ..
  fi
  echo "=> Done"
}
