#!/usr/bin/sh

read test
go:
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
  echo "Буква в верхнем регистре";;
  [3]   ) 
  echo "Цифра";;
esac

