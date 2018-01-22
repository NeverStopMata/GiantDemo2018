#!/bin/bash

docker rm -f pyclient_test1
docker rm -f pyclient_test2
docker rm -f pyclient_test3
docker rm -f pyclient_test4

cp -f _pystress_cfg.json.sample cfg.json
sed -i 's/127.0.0.1/122.11.58.163/g' cfg.json

docker run --name pyclient_test1 -d --restart=always  -v $PWD:/conf py_guiclient
#docker run --name pyclient_test2 -d --restart=always  -v $PWD:/conf py_guiclient
#docker run --name pyclient_test3 -d --restart=always  -v $PWD:/conf py_guiclient
#docker run --name pyclient_test4 -d --restart=always  -v $PWD:/conf py_guiclient
