#!/bin/bash

CUR_DIR=$(pwd)
cd "$REP_DIR" || exit 1

git rm -r --cached .

echo ""
git diff

read -p "Commit Message: " input

git add .
git commit -m "$input"
git push

cd "$CUR_DIR" || exit 1
