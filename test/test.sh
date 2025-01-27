#!/usr/bin/sh

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
esac

