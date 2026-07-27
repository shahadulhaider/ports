#!/usr/bin/env bash
set -euo pipefail

readonly SRC=/opt/demo/listener

# A fresh PID namespace hands out single-digit PIDs, which reads as obviously
# synthetic on screen. Burning ~2400 of them first (≈0.5s) puts the services
# into the 4-digit range a real machine would show.
for _ in $(seq 1 2400); do /bin/true; done

# Each service runs from its own copy of the listener binary: the file name
# becomes the process comm, which is what `ss -tlnp` reports back to the TUI.
spawn() {
  local name=$1
  shift
  install -m 0755 "$SRC" "/usr/local/bin/$name"
  "/usr/local/bin/$name" "$@" &
}

spawn sshd '-tcp6' '[::]:22'
spawn nginx '-tcp4' '0.0.0.0:80,0.0.0.0:443' '-tcp6' '[::]:80'
spawn node '-tcp4' '0.0.0.0:3000,0.0.0.0:5173' '-tcp6' '[::]:3000'
spawn node '-tcp4' '127.0.0.1:4200'
spawn mysqld '-tcp4' '127.0.0.1:3306'
spawn python3 '-tcp4' '127.0.0.1:5000,127.0.0.1:8000'
spawn postgres '-tcp4' '127.0.0.1:5432' '-tcp6' '[::1]:5432'
spawn redis-server '-tcp4' '127.0.0.1:6379'
spawn docker-proxy '-tcp4' '0.0.0.0:8080'
spawn jupyter-lab '-tcp4' '127.0.0.1:8888'
spawn prometheus '-tcp4' '0.0.0.0:9090'
spawn java '-tcp4' '127.0.0.1:9200'
spawn mongod '-tcp4' '127.0.0.1:27017'
spawn memcached '-tcp4' '127.0.0.1:11211'
spawn beam.smp '-tcp4' '0.0.0.0:5672'
spawn dnsmasq '-udp4' '0.0.0.0:53'
spawn avahi-daemon '-udp4' '0.0.0.0:5353'
spawn wireguard-go '-udp4' '0.0.0.0:51820'

sleep 1

spawn psql '-connect' '127.0.0.1:5432' '-dials' '3'
spawn redis-cli '-connect' '127.0.0.1:6379' '-dials' '5'
spawn curl '-connect' '127.0.0.1:3000' '-dials' '2'

sleep 1
clear
exec bash
