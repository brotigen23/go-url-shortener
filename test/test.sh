#!/usr/bin/env bash


secs=3600   # Set interval (duration) in seconds.

SECONDS=0   # Reset $SECONDS; counting of seconds will (re)start from 0(-ish).
while (( SECONDS < secs )); do    # Loop until interval has elapsed.
    curl -i localhost:8080 \
    -X POST \
    -H 'content-type: text/plain' \
    --data "asd" 
    
    curl -i localhost:8080/s3PUAKW4 \
    -X GET 
done

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
