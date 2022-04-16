## 写在前面

- 此教程专为有洁癖的宝宝们准备。不使用任何一键安装脚本。面板党可以退散了！！
- 本人测试环境是 Debian 11 其他的没测试。
## 1.下载源代码
```bash
cd /var/www/
mkdir epusdt
chmod 777 -R /var/www/epusdt
cd epusdt
wget  https://github.com/assimon/epusdt/releases/download/v0.0.2/epusdt_v0.0.2_Linux_x86_64.tar.gz
tar -xzf epusdt_v0.0.2_Linux_x86_64.tar.gz
rm epusdt_v0.0.2_Linux_x86_64.tar.gz
```
## 2.导入Sql
- 创建sql文件
```bash
nano epusdt.sql
```
然后复制下面的
```sql
-- auto-generated definition
use epusdt;
create table orders
(
    id                   int auto_increment
        primary key,
    trade_id             varchar(32)    not null comment 'epusdt订单号',
    order_id             varchar(32)    not null comment '客户交易id',
    block_transaction_id varchar(128)   null comment '区块唯一编号',
    actual_amount        decimal(19, 4) not null comment '订单实际需要支付的金额，保留4位小数',
    amount               decimal(19, 4) not null comment '订单金额，保留4位小数',
    token                varchar(50)    not null comment '所属钱包地址',
    status               int default 1  not null comment '1：等待支付，2：支付成功，3：已过期',
    notify_url           varchar(128)   not null comment '异步回调地址',
    redirect_url         varchar(128)   null comment '同步回调地址',
    callback_num         int default 0  null comment '回调次数',
    callback_confirm     int default 2  null comment '回调是否已确认？ 1是 2否',
    created_at           timestamp      null,
    updated_at           timestamp      null,
    deleted_at           timestamp      null,
    constraint orders_order_id_uindex
        unique (order_id),
    constraint orders_trade_id_uindex
        unique (trade_id)
);

create index orders_block_transaction_id_index
    on orders (block_transaction_id);

-- auto-generated definition
create table wallet_address
(
    id         int auto_increment
        primary key,
    token      varchar(50)   not null comment '钱包token',
    status     int default 1 not null comment '1:启用 2:禁用',
    created_at timestamp     null,
    updated_at timestamp     null,
    deleted_at timestamp     null
)
    comment '钱包表';

create index wallet_address_token_index
    on wallet_address (token);
```
`ctrl+x` 退出，按 `Y`保存 再按回车就好了
- 创建数据库 
```bash
mysql
```
接下来输入命令 
```sql
CREATE DATABASE [这里替换为数据库名] ;
GRANT ALL ON [这里替换为数据库名].* TO '[这里替换为用户名]'@'localhost' IDENTIFIED BY '[这里替换为密码]' WITH GRANT OPTION;
FLUSH PRIVILEGES;
EXIT
```
- 导入sql文件
```bash
mysql -u[用户名] -p[密码] < epusdt.sql 
```
## 3.配置反向代理
```bash
nano /etc/nginx/sites-enabled/epusdt
```
你可以参考以下我的配置文件，注意更改域名。
```bash
server {
   listen 80;
   server_name domain.com;
   return 301 https://domain.com$request_uri;
 }

server {
   listen 443 ssl http2;
   server_name domain.com;
   ssl_certificate  /etc/nginx/sslcert/cert.crt;
   ssl_certificate_key  /etc/nginx/sslcert/key.key; 
   ssl_prefer_server_ciphers on;

   location / {
        proxy_pass http://127.0.0.1:8000;
}
}

```
## 4.赋予Epusdt执行权限
`linux`服务器需要赋予`Epust`执行权限方可启动。            
执行命令```chmod +x epusdt```赋予权限
## 5、配置Epusdt
执行命令
```bash
mv .env.example .env
nano .env
```

```dotenv
app_name=epusdt
#下面配置你的域名，收银台会需要
app_uri=https://upay.dujiaoka.com
#是否开启debug，默认false
app_debug=false
#http服务监听端口
http_listen=:8000
#静态资源文件目录
static_path=/static
#缓存路径
runtime_root_path=/runtime
#日志配置
log_save_path=/logs
log_max_size=32
log_max_age=7
max_backups=3
# mysql配置
mysql_host=127.0.0.1
mysql_port=3306
mysql_user=mysql账号
mysql_passwd=mysql密码
mysql_database=数据库
mysql_table_prefix=
mysql_max_idle_conns=10
mysql_max_open_conns=100
mysql_max_life_time=6
# redis配置
redis_host=127.0.0.1
redis_port=6379
redis_passwd=
redis_db=5
redis_pool_size=5
redis_max_retries=3
redis_idle_timeout=1000
# 消息队列配置
queue_concurrency=10
queue_level_critical=6
queue_level_default=3
queue_level_low=1
#机器人Apitoken
tg_bot_token=
#telegram代理url(大陆地区服务器可使用一台国外服务器做反代tg的url)，如果运行的本来就是境外服务器，则无需填写
tg_proxy=
#管理员userid
tg_manage=
#api接口认证token(用于发起交易的签名认证，请勿外泄)
api_auth_token=
#订单过期时间(单位分钟)
order_expiration_time=10
#强制汇率(设置此参数后每笔交易将按照此汇率计算，例如:6.4)
forced_usdt_rate=
```
⚠️注意：配置文件里面不认识的不要修改，留空即可，不会改又要瞎改，除非你对项目源代码很熟悉很有信心😁
## 6、配置supervisor
为了保证`Epusdt`常驻后台运行，我们需要配置`supervisor`来实现进程监听  
```bash
nano /etc/supervisor/conf.d/epusdt.conf
```
你可以参考以下我的配置文件，注意更改路径。
```conf
[program:epusdt]
process_name=epusdt
directory=/var/www/epusdt
command=/var/www/epusdt/epusdt http start
autostart=true
autorestart=true
user=www-data
numprocs=1
redirect_stderr=true
stdout_logfile=/var/log/supervisor/epusdt.log
```
接下来输入命令
```bash
supervisorctl reread
supervisorctl update
supervisorctl start epusdt
supervisorctl tail epusdt
```
出现下图，即为配置成功
```bash
  _____                     _ _   
 | ____|_ __  _   _ ___  __| | |_ 
 |  _| | '_ \| | | / __|/ _` | __|
 | |___| |_) | |_| \__ \ (_| | |_ 
 |_____| .__/ \__,_|___/\__,_|\__|
       |_|                        
Epusdt version(0.0.2) Powered by assimon https://github.com/assimon/epusdt 
⇨ http server started on [::]:8000
```
## 其他注意事项
- 1.所有`.env`配置文件有了修改后都需要重启supervisor进程 `supervisorctl restart epusdt`
- 2.教程所示的目录均为参考，请勿1:1照抄，根据自己实际情况来
