#!/usr/bin/env bash

svm_try_profile() {
  if [ -z "${1-}" ] || [ ! -f "${1}" ]; then
    return 1
  fi
  echo "${1}"
}

#
# Detect profile file if not specified as environment variable
# (eg: PROFILE=~/.myprofile)
# The echo'ed path is guaranteed to be an existing file
# Otherwise, an empty string is returned
#
svm_detect_profile() {
  if [ "${PROFILE-}" = '/dev/null' ]; then
    # the user has specifically requested NOT to have nvm touch their profile
    return
  fi

  if [ -n "${PROFILE}" ] && [ -f "${PROFILE}" ]; then
    echo "${PROFILE}"
    return
  fi

  local DETECTED_PROFILE
  DETECTED_PROFILE=''

  if [ -n "${BASH_VERSION-}" ]; then
    if [ -f "$HOME/.bashrc" ]; then
      DETECTED_PROFILE="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
      DETECTED_PROFILE="$HOME/.bash_profile"
    fi
  elif [ -n "${ZSH_VERSION-}" ]; then
    DETECTED_PROFILE="$HOME/.zshrc"
  fi

  if [ -n "$DETECTED_PROFILE" ]; then
    echo "$DETECTED_PROFILE"
  fi
}

svm_install_dir() {
  if [ -n "$SVM_DIR" ]; then
    printf %s "${SVM_DIR}"
  else
    printf %s ${HOME}/.svm
  fi
}

svm_clone ()
{
  local INSTALL_DIR
  INSTALL_DIR="$HOME/.svm"

  if [ -f "$INSTALL_DIR/sonic.go" ]; then
    echo "=> svm is already installed in $INSTALL_DIR, trying to update the script"
    cd ~/.svm
    git pull origin master
  else
    echo "=> Cloning svm as script to '$INSTALL_DIR'"
    git clone --quiet --branch feature/copy-skel https://github.com/TeeSrisantitham/sonic.git  ~/.svm
  fi
}

svm_install ()
{
  local SVM_PROFILE
  SVM_PROFILE="$(svm_detect_profile)"
  local PROFILE_INSTALL_DIR
  PROFILE_INSTALL_DIR="$(svm_install_dir)"

  SOURCE_STR="\\nexport SVM_DIR=\"${PROFILE_INSTALL_DIR}\"\\n[ -s \"\$SVM_DIR/scripts/svm.sh\" ] && \\. \"\$SVM_DIR/scripts/svm.sh\""

  if ! command grep -qc '/scripts/svm.sh' "$SVM_PROFILE"; then
    echo "=> Appending svm source string to $SVM_PROFILE"
    command printf "${SOURCE_STR}" >> "$SVM_PROFILE"
  else
    echo "=> svm is already append in $SVM_PROFILE"
  fi
}

svm()
{
  svm_clone
  svm_install
}

svm
