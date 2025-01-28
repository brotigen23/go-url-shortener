#!/usr/bin/env bash

read test
case "$test" in
  [1]   ) 
  echo "Default test"
  
    curl -i localhost:8080 \
    -X POST \
    -H 'content-type: text/plain' \
    --data "asd"  \

    curl -i localhost:8080 \
    -X POST \
    -H 'content-type: application/json' \
    --data '{"url":"ya.ru"}'  \

    curl -i localhost:8080/api/shorten/batch \
    -X POST \
    -H 'content-type: application/json' \
    --data '
[
  {
    "correlation_id":"1",
    "original_url":"example.com"
  }, 
  {   
    "correlation_id":"2",
    "original_url":"example1.com"
  }
]'
;;
  [2]   )   
echo "get user's urls"
curl -i localhost:8080/api/user/urls \
-X GET \
-H 'content-type: application/json' \
--cookie 'JWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDE2ODMxOTAsIlVzZXJuYW1lIjoiMndnSTl3UFRMMmJ6Z3pzUiJ9.gb184S3gyHcNjrxAr0yImWRUWX_6hSRnH0zStpTxnXE'\
  ;;
  [3]   ) 
curl -i localhost:8080 \
-X POST \
-H 'content-type: text/plain' \
--data "asd"  \
  ;;
  [4]   ) 
curl -i localhost:8080/api/user/urls \
-X DELETE \
-H 'content-type: application/json' \
--cookie 'JWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDE2OTg4OTEsIlVzZXJuYW1lIjoiU2I0ZlRVRUltcklHNndSVyJ9.u2Jxoqwu6lxdISgQoapg4hv0zOlJF4x7iPkP7zZu8LY' \
--data '
[
  "lmMJw0Ho", "HET0DTqZ"
]'
;;
esac
