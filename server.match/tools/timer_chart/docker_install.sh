#!/bin/bash

docker rm -f py_timer_chart
docker run --name py_timer_chart -d --restart=always py_timer_chart --ip=122.11.58.163
