#!/bin/bash

echo GITHUB_REPO_NAME
read repoName
echo GITHUB_REPO_ORGANIZATION_PATH
read gitPath

rm -rf ~/.sonic
git clone https://github.com/openware/sonic.git ~/.sonic

mkdir -p $repoName
cp -R ~/.sonic/skel/ ./$repoName

cd $repoName
git init
git remote add origin $gitPath
