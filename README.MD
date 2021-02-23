# 微服务网关中间件

注意事项：
1. 仅能运行在Linux服务器上，使用容器环境。
1. 统一内部HTTP服务端口为80，GRPC端口为8000。
1. 网关中间件监听端口为：HTTP/10080, GRPC/18000。
1. 网关中间件通过`SideCar`模式与业务容器运行，需要通过`iptables`定义转发规则，将所有内网请求转发到网关中间件：
    ```
    iptables -t nat -A OUTPUT -p tcp -m iprange --dst-range 10.0.0.0-10.255.255.255 --dport 80 -j DNAT --to-destination 127.0.0.1:10080
    iptables -t nat -A OUTPUT -p tcp -m iprange --dst-range 172.16.0.0-172.31.255.255 --dport 80 -j DNAT --to-destination 127.0.0.1:10080
    iptables -t nat -A OUTPUT -p tcp -m iprange --dst-range 192.168.0.0-192.168.255.255 --dport 80 -j DNAT --to-destination 127.0.0.1:10080

    iptables -t nat -A OUTPUT -p tcp -m iprange --dst-range 10.0.0.0-10.255.255.255 --dport 8000 -j DNAT --to-destination 127.0.0.1:18000
    iptables -t nat -A OUTPUT -p tcp -m iprange --dst-range 172.16.0.0-172.31.255.255 --dport 8000 -j DNAT --to-destination 127.0.0.1:18000
    iptables -t nat -A OUTPUT -p tcp -m iprange --dst-range 192.168.0.0-192.168.255.255 --dport 8000 -j DNAT --to-destination 127.0.0.1:18000
    ```
    在容器中可以通过以下命令查看`iptables`转发规则：
    ```
    iptables -t nat --list
    ```
1. `SideCar`模式运行时，务必保证`Mesh`组件优先于所有容器运行。


# 注意事项

1. 目前该项目只实现了基本的转发能力，仅供学习参考。
2. 后续需要增加比较重要的负载均衡能力，依靠本身的GRPC负载均衡实现。
3. `SideCar`模式的一些常见痛点解决请参考：https://banzaicloud.com/blog/k8s-sidecars/
