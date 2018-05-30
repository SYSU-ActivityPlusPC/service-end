# SYSU ActivityPlus service for PC

## 文件结构说明

- controller/

  负责对请求进行处理和转发，主要的业务逻辑都是在这之中完成的

- model/

  对数据进行管理

- router/

  路由函数的存放地址

## 环境变量

- STATIC_DIR

  静态文件夹位置

- PORT

  端口号

- DATABASE_ADDRESS

  数据库地址

- DATABASE_PORT

  数据库端口

- ADMIN_MAIL_PASS

  管理员账户的邮箱密码

- MQ_ADDRESS

  消息队列的地址，默认为`localhost`

- MQ_PORT

  消息队列端口，默认为`5672`

- MQ_USER

  消息队列用户名

- MQ_PASSWORD

  消息队列密码