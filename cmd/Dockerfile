# 基础镜像
FROM ubantu:20.04
# 把编译后的打包进来这个镜像，放到工作目录 /app。你随便换
COPY ebook /app/ebook
WORKDIR /app
# CMD 是执行命令
# 最佳
ENTRYPOINT ["/app/ebook"]