#!/bin/bash

usage(){
  echo "$1 argument is mantatory"
  echo "TARGET get only be server, client, both, all"
  echo "Example: $0 server"
  echo "Example: $0 client"
  echo "Example: $0 both"
  echo "Example: $0 all"
}

if [ -n "$1" ]
then
  export TARGET="$1"
else
  usage TARGET
  exit 1
fi
shift

TMP_FILE=/tmp/vprovison$$
vagrant status &> ${TMP_FILE}

SERVERS=$(cat ${TMP_FILE} | grep server | awk '{print $1}')
CLIENTS=$(cat ${TMP_FILE} | grep client | awk '{print $1}')

if [[ "${TARGET}" == "server" ]]
then
  for i in ${SERVERS}; do vagrant provision "server$i" --provision-with scalezilla; done
fi
shift

if [[ "${TARGET}" == "client" ]]
then
  for i in ${CLIENTS}; do vagrant provision $i --provision-with scalezilla; done
fi
shift

if [[ "${TARGET}" == "both" ]]
then
  for i in ${SERVERS}; do vagrant provision $i --provision-with scalezilla; done
  for i in ${CLIENTS}; do vagrant provision $i --provision-with scalezilla; done
fi

if [[ "${TARGET}" == "all" ]]
then
  vagrant provision dev --provision-with scalezilla
  for i in ${SERVERS}; do vagrant provision $i --provision-with scalezilla; done
  for i in ${CLIENTS}; do vagrant provision $i --provision-with scalezilla; done
fi

rm -f ${TMP_FILE}

