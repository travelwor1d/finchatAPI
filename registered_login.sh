#!/bin/bash

set -e

if [ -z $API_KEY ]; then
  echo 'Please provide API_KEY'
  exit 1
fi;

curl -X POST -H 'Content-Type: application/json' -d '
{
  "email": '"'$1'"',
  "password": '"'$2'"',
  "returnSecureToken": true
}
' "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=$API_KEY"