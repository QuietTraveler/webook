# 使用基础镜像
FROM ubuntu:20.04

# 创建工作目录 /app
WORKDIR /app

# 把编译后的文件放入 /app 目录
COPY webook /app/webook

# 确保 ENTRYPOINT 运行的是 /app/webook
ENTRYPOINT ["/app/webook"]
