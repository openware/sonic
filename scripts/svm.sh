#!/usr/bin/env bash

svm() {
  if [ $# -lt 1 ]; then
    svm --help
    return
  fi

  local COMMAND
  COMMAND="${1}"

  case ${COMMAND} in
  '-h' | 'help' | '--help')
    echo 'Usage:'
    echo '  svm --help                                  Show this message'
    echo '  svm create <project name>                   Install sonic'
    return
    ;;
  'create' | 'i')
    create "${2}"
    ;;
  *)
    hander_err "Command ${COMMAND} is not found"
    ;;
  esac
}

hander_err() {
  echo "=> Error: ${1}"
  exit 1
}

create() {
  local GITPATH
  GITPATH="${1}"
  local DIR
  local ADDR

  IFS='/'
  read -A ADDR <<<"$GITPATH"
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
