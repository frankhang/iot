#!/bin/bash

#for linux

version=1.1

if [ $# -gt 2 ]; then
  echo "Usage: $0 [service] [-r]"
  echo "  [service]: prom/grafana/vehicle/client/server, server means prom + grafana + vehicle. Default is blank, means start them all"
  echo "  [-r]: remove image before running"
  exit 0
fi

if [ "$1" = "" -o "$1" = "-r" ]; then
  grafana=true
  prom=true
  vehicle=true
  client=true
fi

arg=ops

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

if [ "$1" = "-r" -o "$2" = "-r" ]; then
  remove=true
fi

if [ $grafana ]; then
  echo -e
  echo "#### staring grafana ####"
  image=frankhang/grafana:$version
  if [ $remove ]; then
    docker image rm -f $image
  fi
  docker stop grafana
  docker rm grafana
  mkdir /grafanadata
  chmod a+rwx /grafanadata
  docker run --name grafana -d --network=host --add-host=host.docker.internal:127.0.0.1 -v /grafanadata:/var/lib/grafana $image
fi

if [ $prom ]; then
  echo -e
  echo "#### staring prom ####"
  image=frankhang/prom-$arg:$version
  if [ $remove ]; then
    docker image rm -f $image
  fi
  docker stop prom
  docker rm prom

  mkdir /promdata
  chmod a+rwx /promdata
  docker run --name prom -d --network=host --add-host=host.docker.internal:127.0.0.1 -v /promdata:/etc/prometheus/data $image

fi

if [ $vehicle ]; then
  echo -e
  echo "#### staring vehicle ####"
  image=frankhang/iot-vehicle:$version
  if [ $remove ]; then
    docker image rm -f $image
  fi
  docker stop vehicle
  docker rm vehicle
  docker run --name vehicle -d --network=host --add-host=host.docker.internal:127.0.0.1 $image -L=debug

fi

if [ $client ]; then
  echo -e
  echo "#### starting client1  ####"
  image=frankhang/client:$version
  if [ $remove ]; then
    docker image rm -f $image
  fi
  docker stop client1
  docker rm client1
  docker stop client2
  docker rm client2
  docker run --name client1 -d --network=host --add-host=host.docker.internal:127.0.0.1 $image
  docker run --name client2 -d --network=host --add-host=host.docker.internal:127.0.0.1 $image
fi
