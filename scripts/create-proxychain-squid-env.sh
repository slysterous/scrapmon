#!/bin/bash

rm ./proxy-list.txt
curl "https://api.proxyscrape.com/?request=getproxies&proxytype=http&timeout=10000&country=all&ssl=no&anonymity=all" -o ./proxy-list.txt


rm ./proxychains.conf
touch ./proxychains.conf
# chown ./proxychains.conf 777

# echo "dynamic_chain" >> ./proxychains.conf
# echo "[ProxyList]" >> ./proxychains.conf
# echo "PROXYCHAIN=\"yes\"" >> ./squid-proxy-list.env
# echo "PROXYCHAIN_TYPE=\"dynamic_chain\"" >> ./squid-proxy-list.env

input="./proxy-list.txt"

index=1

cat $input | while read line 
do
  IFS=: read -r ipVar portVar <<< "$line"

  counter=${index}
  protocol='http '
  ip="${ipVar}"
  port="${portVar::-1}"
  echo -e "${protocol}${ip} ${port}" >> ./proxychains.conf

  index=$(($index+1))
   
done