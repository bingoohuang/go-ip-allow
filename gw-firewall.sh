#!/bin/bash

IPTABLES=/sbin/iptables
IPTABLES_SAVE=/sbin/iptables-save

$IPTABLES_SAVE > /usr/local/etc/iptables.last

echo '配置默策略ACCEPT'
$IPTABLES -P INPUT ACCEPT
$IPTABLES -F INPUT

echo '清空nat表'
$IPTABLES -t nat -F
$IPTABLES -t nat -X
$IPTABLES -t nat -Z

$IPTABLES -t mangle -F
$IPTABLES -t mangle -X
$IPTABLES -t mangle -Z

echo '删除admin链'
$IPTABLES -F admin
$IPTABLES -X admin

echo '清空filter表'
$IPTABLES -t filter -F
$IPTABLES -t filter -X
$IPTABLES -t filter -Z

echo '允许本机访问'
$IPTABLES -A INPUT -s 127.0.0.1 -j ACCEPT
$IPTABLES -A INPUT -s 127.0.0.1 -d 127.0.0.1 -j ACCEPT
$IPTABLES -A OUTPUT -j ACCEPT

echo '允许指定IP或者IP段访问，固定开放的IP'
$IPTABLES -A INPUT --src 192.168.0.0/16 -j ACCEPT
$IPTABLES -A INPUT --src 121.40.25.37 -j ACCEPT  #svm
$IPTABLES -A INPUT --src 121.40.27.88 -j ACCEPT  #官网

ALLOW_IPS=121.237.61.142,121.237.61.152

echo '创建新的链 admin，临时开放的IP'
$IPTABLES -N admin
$IPTABLES -A admin --src $ALLOW_IPS -j ACCEPT
$IPTABLES -A admin -j DROP
$IPTABLES -A INPUT -m tcp -p tcp -m multiport --dports 22,13306 -j admin # allow 22 and redis
$IPTABLES -A INPUT -m tcp -p tcp -j admin


echo '配置nat表 DNAT与SNAT'
$IPTABLES -I INPUT -p tcp --dport 8080  -j ACCEPT  # CI
$IPTABLES -t nat -I PREROUTING -p tcp --dport 1122 -j DNAT --to 192.168.1.1:22
$IPTABLES -t nat -I PREROUTING -p tcp --dport 1222 -j DNAT --to 192.168.1.2:22
$IPTABLES -t nat -I POSTROUTING -s 192.168.0.0/16 -j SNAT --to-source 192.168.0.1
# $IPTABLES -I INPUT -m state --state NEW -m tcp -p tcp --dport 8000 -m mark --mark 1 -j ACCEPT
# $IPTABLES -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8000
# $IPTABLES -t mangle -A PREROUTING -p tcp --dport 80 -j MARK --set-mark 1
$IPTABLES -t nat -A PREROUTING -p tcp --dport 1000 -j REDIRECT --to-port 8000


echo '配置filter表'
$IPTABLES -I INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT
$IPTABLES -A INPUT -p udp -j DROP
$IPTABLES -A INPUT -p tcp --syn -j DROP

/etc/init.d/iptables save
