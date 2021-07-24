#!/bin/bash

#for linux


#!/bin/bash

#for mac

version=1.2



if [ $# = 0 ]; then
  grafana=true
  prom=true
  vehicle=true
  client=true
fi

if [ "$1" = "server" ]; then
  grafana=true
  prom=true
  vehicle=true
fi

if [ "$1" = "grafana" ]; then
  grafana=true
fi

if [ "$1" = "prom" ]; then
  prom=true
fi

if [ "$1" = "vehicle" ]; then
  vehicle=true
fi

if [ "$1" = "client" ]; then
  client=true
fi



if [ $grafana ]; then
  echo -e
  echo "#### stoping grafana ####"
  docker stop grafana
  docker rm grafana
fi

if [ $prom ]; then
  echo -e
  echo "#### stoping prom ####"
  docker stop prom
  docker rm prom
fi

if [ $vehicle ]; then
  echo -e
  echo "#### stoping vehicle ####"
  docker stop vehicle
  docker rm vehicle

fi

if [ $client ]; then
  echo -e
  echo "#### stoping client  ####"
  docker stop client1
  docker rm client1
  docker stop client2
  docker rm client2
fi





