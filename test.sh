#!/bin/bash

curl -X POST http://localhost:8081/csv \
  -F "email=coolorvibes@gmail.com" \
  -F "file=@txns.csv"