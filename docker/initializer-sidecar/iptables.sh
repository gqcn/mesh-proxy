#!/bin/bash

# 注意如果是用户UID:1880发起的请求不会走iptables规则。

# HTTP服务转发，仅支持内网请求
iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner 1880 -m iprange --dst-range 10.0.0.0-10.255.255.255 --dport 80 -j DNAT --to-destination 127.0.0.1:10080
iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner 1880 -m iprange --dst-range 172.16.0.0-172.31.255.255 --dport 80 -j DNAT --to-destination 127.0.0.1:10080
iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner 1880 -m iprange --dst-range 192.168.0.0-192.168.255.255 --dport 80 -j DNAT --to-destination 127.0.0.1:10080

# GRPC服务转发，仅支持内网请求
iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner 1880 -m iprange --dst-range 10.0.0.0-10.255.255.255 --dport 8000 -j DNAT --to-destination 127.0.0.1:18000
iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner 1880 -m iprange --dst-range 172.16.0.0-172.31.255.255 --dport 8000 -j DNAT --to-destination 127.0.0.1:18000
iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner 1880 -m iprange --dst-range 192.168.0.0-192.168.255.255 --dport 8000 -j DNAT --to-destination 127.0.0.1:18000

# 查看转发列表
iptables -t nat --list