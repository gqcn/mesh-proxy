# loads/mesh-proxy:test
# loads/mesh-proxy:latest

# gf pack config boot/data-packed.go -n boot
# gf docker -t loads/mesh-proxy:prod -p
# gf docker -t loads/mesh-proxy:test -p
# gf docker -t loads/mesh-proxy:dev -p
FROM loads/alpine:3.8

LABEL maintainer="john@johng.cn"

###############################################################################
#                                INSTALLATION
###############################################################################

# 环境变量设置
ENV APP_NAME mesh-proxy
ENV APP_ROOT /var/www
ENV APP_PATH $APP_ROOT/$APP_NAME
ENV LOG_ROOT /var/log/www
ENV LOG_PATH /var/log/www/$APP_NAME

# 关闭自动日志搜集，改为手动
ENV LOG_TAIL 0

# 创建必要的运行用户, UID固定
RUN adduser proxy -D -u 1880 -s /bin/bash

# 执行入口文件添加
ADD ./bin/linux_amd64/main $APP_PATH/
RUN chmod 777 $APP_PATH -R
ADD ./docker/dockerfiles/*.sh /bin/
RUN chmod +x /bin/*.sh

###############################################################################
#                                   START
###############################################################################

USER proxy
CMD  proxy.sh