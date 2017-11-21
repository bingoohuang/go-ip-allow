#!/usr/bin/env bash

if [ $# == 0 ]; then
        echo "请输入两个参数，参数1为环境名列表，例如：hb,hd,dev(以逗号分隔); 参数2表示ip列表(以逗号分隔)"
        exit 1
fi


sed -i '' "s/^ALLOW_IPS.*/ALLOW_IPS=$1/"  "gw-firewall.sh"

