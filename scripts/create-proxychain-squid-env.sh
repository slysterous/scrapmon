#!/bin/bash


curl "https://api.proxyscrape.com/?request=getproxies&proxytype=http&timeout=10000&country=all&ssl=no&anonymity=all" -o ./proxy-list.txt

rm ./squid-proxu-list.env
touch ./squid-proxy-list.env

echo "VISIBLE_HOSTNAME=docker-squid4-print-scrape" >> ./squid-proxy-list.env
echo "PROXYCHAIN=yes" >> ./squid-proxy-list.env
echo "PROXYCHAIN_TYPE=dynamic_chain" >> ./squid-proxy-list.env

input="./proxy-list.txt"

index=1

while IFS= read -r -uN line
do
  echo $line
  IFS=', ' read -r -a array <<< "$line"
  # Storing as array into myarray
  args=$(echo $line |tr -d '\r\n' | tr ":" "\n" )

  envline=PROXYCHAIN_PROXY
  envline+="${index}=\"http "
  envline+="${args}\""
  echo ${envline} >> ./squid-proxy-list.env

  index=$(($index+1))
  if [ $index -gt 999 ]
  then 
    break
  fi

done N< $input