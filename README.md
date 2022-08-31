## Epusdt (Easy Payment Usdt)
<p align="center">
<img src="wiki/img/usdtlogo.png">
</p>
<p align="center">
<a href="https://www.gnu.org/licenses/gpl-3.0.html"><img src="https://img.shields.io/badge/license-GPLV3-blue" alt="license GPLV3"></a>
<a href="https://golang.org"><img src="https://img.shields.io/badge/Golang-1.16-red" alt="Go version 1.16"></a>
<a href="https://echo.labstack.com"><img src="https://img.shields.io/badge/Echo Framework-v4-blue" alt="Echo Framework v4"></a>
<a href="https://github.com/tucnak/telebot"><img src="https://img.shields.io/badge/Telebot Framework-v3-lightgrey" alt="Telebot Framework-v3"></a>
<a href="https://github.com/assimon/epusdt/releases/tag/v0.0.1"><img src="https://img.shields.io/badge/version-v0.0.1-green" alt="version v0.0.1"></a>
</p>


## 项目简介
`Epusdt`（全称：Easy Payment Usdt）是一个由`Go语言`编写的私有化部署`Usdt`支付中间件(`Trc20网络`)     
站长或开发者可通过`Epusdt`提供的`http api`集成至您的任何系统，无需过多的配置，仅仅依赖`mysql`和`redis`      
可实现USDT的在线支付和消息回调，这一切在优雅和顷刻间完成！🎉        
私有化搭建使得无需额外的手续费和签约费用，Usdt代币直接进入您的钱包💰      
`Epusdt` 遵守 [GPLv3](https://www.gnu.org/licenses/gpl-3.0.html) 开源协议!

## 项目特点
- 支持私有化部署，无需担心钱包被篡改和吞单😁
- `Go语言`跨平台实现，支持x86和arm芯片架构的win/linux设备
- 多钱包地址轮询，提高订单并发率
- 异步队列响应，优雅及高性能
- 无需额外环境配置，仅运行一个编译后二进制文件即可使用
- 支持`http api`，其他系统亦可接入
- `Telegram`机器人接入，便捷使用和支付消息快速通知

## 项目结构
```
Epusdt
    ├── plugins ---> (已集成的插件库，例如dujiaoka)
    ├── src ---> (项目核心目录）
    ├── sdk ---> (接入SDK)
    ├── sql ---> (安装sql文件或更新sql文件)
    └── wiki ---> (知识库)
```

## 教程：
- 宝塔运行`epusdt`教程👉🏻[宝塔运行epusdt](wiki/BT_RUN.md)
- 不好意思我有洁癖，手动运行`epusdt`教程👉🏻[手动运行epusdt](wiki/manual_RUN.md)
- 开发者接入`epusdt`文档👉🏻[开发者接入epusdt](wiki/API.md)
- HTML+PHP极速运行`epusdt`教程👉🏻[使用PHPAPI-for-epusdt极速接入epusdt](https://github.com/BlueSkyXN/PHPAPI-for-epusdt)

## 已适配系统插件
- 独角数卡[插件地址](plugins/dujiaoka)

## 🔥推荐服务器 
- （美国免备案vps，配置2核2G仅需`20.98$`≈`145RMB`一年/支持支付宝付款）[👉🏻点我直达](https://my.racknerd.com/aff.php?aff=2745&pid=681)

## 加入交流/意见反馈
- `Epusdt`频道[https://t.me/epusdt](https://t.me/epusdt)
- `Epusdt`交流群组[https://t.me/epusdt_group](https://t.me/epusdt_group)

## 设计实现
`Epusdt`的实现方式与其他项目原理类似，都是通过监听`trc20`网络的api或节点，      
监听钱包地址`usdt`代币入账事件，通过`金额差异`和`时效性`来判定交易归属信息，     
可参考下方`流程图`
```
简单的原理：
1.客户需要支付20.05usdt
2.服务器有一个hash表存储钱包地址对应的待支付金额 例如:address_1 : 20.05
3.发起支付的时候，我们可以判定钱包address_1的20.05金额是否被占用，如果没有被占用那么可以直接返回这个钱包地址和金额给客户，告知客户需按规定金额20.05准确支付，少一分都不行。且将钱包地址和金额 address_1:20.05锁起来，有效期10分钟。
4.如果订单并发下，又有一个20.05元需要支付，但是在第3步的时候上一个客户已经锁定了该金额，还在等待支付中...，那么我们将待支付金额加上0.0001，再次尝试判断address_1:20.0501金额是否被占用？如果没有则重复第三步，如果还是被占用就继续累加尝试，直到加了100次后都失败
5.新开一个线程去监听所有钱包的USDT入账事件，网上有公开的api或rpc节点。如果发现有入账金额与待支付的金额相等。则判断该笔订单支付成功！
```
### 流程图：
![Implementation principle](wiki/img/implementation_principle.jpg)

## 打赏
如果该项目对您有所帮助，希望可以请我喝一杯咖啡☕️
```
Usdt(trc20)打赏地址: TNEns8t9jbWENbStkQdVQtHMGpbsYsQjZK
```
<img src="wiki/img/usdt_thanks.jpeg" width = "300" height = "400" alt="usdt扫码打赏"/>




## 声明
`Epusdt`为开源的产品，仅用于学习交流使用！       
不可用于任何违反中华人民共和国(含台湾省)或使用者所在地区法律法规的用途。           
因为作者即本人仅完成代码的开发和开源活动(开源即任何人都可以下载使用或修改分发)，从未参与用户的任何运营和盈利活动。       
且不知晓用户后续将程序源代码用于何种用途，故用户使用过程中所带来的任何法律责任即由用户自己承担。            
```
！！！Warning！！！
项目中所涉及区块链代币均为学习用途，作者并不赞成区块链所繁衍出代币的金融属性
亦不鼓励和支持任何"挖矿"，"炒币"，"虚拟币ICO"等非法行为
虚拟币市场行为不受监管要求和控制，投资交易需谨慎，仅供学习区块链知识
```
