# Ontology平台TPS测试报告
## 测试目的
TPS (Transactions Per Second) 即*每秒处理交易数*，是衡量区块链性能的一个重要指标。本次测试旨在基于特定的测试环境下量化地评估本体在不同的TPS压力下的交易处理能力，以及其处理能力随TPS压力变化而变化的趋势。本次测试将分别测试本体在solo (单节点) 模式和dBFT共识模式下的TPS。

## 测试环境
### solo模式
发送交易机器：
- 硬件：Intel(R) Xeon(R) CPU E5-2430 | 8GB RAM | 120GB　SSD
- 软件：CentOS Linux release 7.6.1810 (Core) | JMeter 5.0
  
运行节点机器：
- 硬件：一台 | Intel(R) Xeon(R) CPU E5-2680 v2 | 64GB RAM | 500GB SSD
- 软件：CentOS Linux release 7.6.1810 (Core) | go1.11 linux/amd64 | ont v1.6.0-beta-32-g8f2c9e9

网络带宽：
- 100Mbit/s，但测试过程中实测只有5Mbit/s
### dBFT模式
发送交易机器：
- 硬件：Intel(R) Xeon(R) CPU E5-2430 | 8GB RAM | 120GB　SSD
- 软件：CentOS Linux release 7.6.1810 (Core) | JMeter 5.0
  
运行节点机器：
- 硬件：四台 | Intel(R) Xeon(R) CPU E5-2430 | 8GB RAM | 120GB　SSD
- 软件：CentOS Linux release 7.6.1810 (Core) | go1.11 linux/amd64 | ont v1.6.0-beta-32-g8f2c9e9

网络带宽：
- 100Mbit/s，但测试过程中实测只有5Mbit/s
## 测试方法
### 部署节点
1. 在运行节点机器上安装好[开发环境](https://github.com/ontio/ontology#build-development-environment)后，从Github上下载ontology源码到$GOPATH/src/github.com/ontio目录下：
   > $ git clone https://github.com/ontio/ontology.git  
   或者  
   > $ go get github.com/ontio/ontology

2. 然后需要修改三处代码以突破输入TPS限制：
* ontology/main.go文件的第283行后面插入以下内容：
   ```go
    stlValidator3, _ := stateless.NewValidator("stateless_validator3")
	stlValidator3.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator4, _ := stateless.NewValidator("stateless_validator4")
	stlValidator4.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator5, _ := stateless.NewValidator("stateless_validator5")
	stlValidator5.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator6, _ := stateless.NewValidator("stateless_validator6")
	stlValidator6.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator7, _ := stateless.NewValidator("stateless_validator7")
	stlValidator7.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator8, _ := stateless.NewValidator("stateless_validator8")
	stlValidator8.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator9, _ := stateless.NewValidator("stateless_validator9")
	stlValidator9.Register(txPoolServer.GetPID(tc.VerifyRspActor))
   ```
   修改后的文件内容部分如下所示：
   ```go
   ...
    stlValidator, _ := stateless.NewValidator("stateless_validator")
	stlValidator.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator2, _ := stateless.NewValidator("stateless_validator2")
	stlValidator2.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator3, _ := stateless.NewValidator("stateless_validator3")
	stlValidator3.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator4, _ := stateless.NewValidator("stateless_validator4")
	stlValidator4.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator5, _ := stateless.NewValidator("stateless_validator5")
	stlValidator5.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator6, _ := stateless.NewValidator("stateless_validator6")
	stlValidator6.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator7, _ := stateless.NewValidator("stateless_validator7")
	stlValidator7.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator8, _ := stateless.NewValidator("stateless_validator8")
	stlValidator8.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator9, _ := stateless.NewValidator("stateless_validator9")
	stlValidator9.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stfValidator, _ := stateful.NewValidator("stateful_validator")
    stfValidator.Register(txPoolServer.GetPID(tc.VerifyRspActor))
   ...
   ```
* ontology/p2pserver/common/p2p_common.go文件的55行修改为以下内容：
   ```go
   MAX_MSG_LEN = 60 * 1024 * 1024    //the maximum message length
   ```
* ontology/txnpool/common/txnpool_common.go文件的28到30行修改为以下内容：
   ```go
   MAX_CAPACITY     = 10014000       // The tx pool's capacity that holds the verified txs
   MAX_PENDING_TXN  = 409600 * 10    // The max length of pending txs
   MAX_WORKER_NUM   = 20             // The max concurrent workers
   ```
3. 用glide安装/更新依赖项：
   > $ cd \$GOPATH/src/github.com/ontio/ontology   
   > $ glide install
   > $ glide update
4. 用make指令生成可执行文件:
   > $ make all
5. 每个节点下面添加一个wallet.dat文件：
   > ./ontology account add --default
6. solo模式下的启动命令为 (solo模式下最短出块时间为2秒)：
   > ./ontology --testmode --testmode-gen-block-time 2

   dBFT模式下需修改[myDBFT.json配置文件](https://github.com/joeqian10/ont_TPS_testing/blob/master/OntConfig/myDBFT.json)中的"SeedList"和"Bookkeepers"，并将其放到ontology目录下 (即和ontology/main.go处于同一个文件夹)，启动命令为：
   > ./ontology --config myDBFT.json --enable-consensus --rest --gaslimit 20000 --gasprice 0

### JMeter配置
1. 在发送交易机器上安装Java Development Kit：
   > yum install openjdk-8-jdk
2. 下载JMeter：
   > wget -c http://ftp.ps.pl/pub/apache//jmeter/binaries/apache-jmeter-5.0.zip
3. 解压：
   > unzip apache-jmeter-5.0.zip
4. 添加.jar文件到JMeter：
   把[这两个.jar文件](https://github.com/joeqian10/ont_TPS_testing/tree/master/jmeter-support-jar)添加到jmeter-install-path/bin/目录下
5. 将JMeter的[.jmx配置文件](https://github.com/joeqian10/ont_TPS_testing/blob/master/JMeterConfig)添加到 jmeter-install-path/bin/目录下，将文件中的服务器IP地址改为本体网络所使用的IP地址，可以根据需要修改线程数和固定吞吐量定时器(ConstantThroughputTimer)的发送速率
6. 运行JMeter：
   > ./jmeter -n -t solo.jmx   
   
   或者  
   > ./jmeter -n -t dbft.jmx 

### 构造交易
1. 在发送交易机器上安装好[开发环境](https://github.com/ontio/ontology#build-development-environment)后，从Github上下载ontology和ontology-go-sdk源码到$GOPATH/src/github.com/ontio目录下：
   > $ git clone https://github.com/ontio/ontology.git  
   > $ git clone https://github.com/ontio/ontology-go-sdk.git   
   
   或者  
   > $ go get github.com/ontio/ontology   
   > $ go get github.com/ontio/ontology-go-sdk
2. solo模式下将节点的wallet.dat拷贝到jmeter-install-path/bin/目录下，dBFT模式下将1号节点(.jmx配置文件中的server1)的wallet.dat拷贝到jmeter-install-path/bin/目录下，并命名为wallet-admin.dat
3. 新生成一个wallet.dat，将其拷贝到jmeter-install-path/bin/目录下，并命名为wallet-account.dat
4. 修改[构造交易源码](https://github.com/joeqian10/ont_TPS_testing/tree/master/GenerateTransactions)，将第23行中的ont地址替换为wallet-admin.dat中的ont地址，将第30行中的ont地址替换为wallet-account.dat中的ont地址
5. 编译修改后的代码：
   > go build -o gen-txs main.go
6. 将生成的文件gen-txs拷贝到jmeter-install-path/bin/目录下
7. 运行gen-txs构造交易，默认生成1000条交易，可加入一个运行参数n，则生成n*1000条交易，最大不超过2^32条交易，比如：
   > gen-txs 5 >transfer.txt      
   
   会在transfer.txt中包含5000条交易的hash和表示为HexString的序列化后的交易数据

### 多签地址转账
solo模式下节点账户中已有10亿ONT，dBFT模式下可以使用[多签地址转账工具](https://github.com/joeqian10/ont_TPS_testing/tree/master/TransferFromMultiSignAddressTool)将10亿ONT转到1号节点(.jmx配置文件中的server1)的wallet.dat的ont地址中，具体步骤如下：
1. 将/TransferFromMultiSignAddressTool/wallets/path1/wallet.dat替换为1号节点的wallet.dat
2. 将/TransferFromMultiSignAddressTool/wallets/path2/wallet.dat替换为2号节点的wallet.dat
3. 将/TransferFromMultiSignAddressTool/wallets/path3/wallet.dat替换为3号节点的wallet.dat
4. 将/TransferFromMultiSignAddressTool/wallets/path4/wallet.dat替换为4号节点的wallet.dat
5. 修改[多签地址转账源码](https://github.com/joeqian10/ont_TPS_testing/tree/master/TransferFromMultiSignAddressTool/main.go)，将第20行中ont地址替换为1号节点的wallet.dat的ont地址
6. 编译修改后的代码：
   > go build main.go
7. 启动ontology测试网络后运行代码即可：
   > ./main

### 测试步骤
1. [构造交易](#构造交易)，这批交易的总数量可根据待测的输入TPS而定(比如，输入TPS为5时，可构造1000笔交易，大约200秒发送完毕；输入TPS为100时，可构造30000笔，大约300秒发送完毕；等等)，每笔转账的数额为1；
2. [部署节点](#部署节点)并启动ontology测试网络，dBFT模式下可以使用[多签地址转账](#多签地址转账)将10亿ONT转到一个节点账户中；
3. 查询转账交易的接收账户的*初始余额*(可通过节点的RPC接口发送查询请求)，并记录；
4. [配置JMeter](#JMeter配置)并将步骤1中的这批交易发送到测试网络上，交易数量较多时可配置多个线程发送；
5. 查看节点日志，观察每个块中的落账交易数量，出现第一个非空块时记录时间作为*起始时间*，发送完毕后，连续出现三个以上的空块时，可认为交易已经处理完毕，取最后一个非空块的时间作为落账*结束时间*；
6. 查询转账交易的接收账户的*结束余额*，并记录；
7. 在这一输入TPS下的输出TPS = ( *结束余额* - *初始余额* ) / ( *结束时间* - *起始时间* )；
8. 在新的输入TPS下，重复步骤1到7。

## 测试结果
### solo模式
![solo_TPS](https://github.com/joeqian10/ont_TPS_testing/blob/master/testingResult/solo_TPS.png)

### dBFT模式
![dbft_TPS](https://github.com/joeqian10/ont_TPS_testing/blob/master/testingResult/dbft_TPS.png)

### 总结
从测试结果可以看出，输出TPS随输出TPS的变化趋势为：
1. 在较低的输入TPS下，输出TPS和输入TPS基本保持一致；
2. 随着输入TPS的增加，两种模式都在一定的TPS下达到峰值，并在输入TPS的继续增大下保持稳定(solo模式下峰值输出TPS为600左右，dBFT模式下为240左右)。   

本测试结果只在此特定的测试环境下具有参考价值。

## 参考
1. [Ontology VBFT共识算法性能测试报告](https://my.oschina.net/uchihamadara/blog/1828848)，非常感谢此文作者对本测试提供的帮助。
2. [NEO平台TPS压力测试报告](http://47.102.110.9:8090/pages/viewpage.action?pageId=65818)
3. [本体开发者文档](https://dev-docs.ont.io/#/docs-cn/)
4. [ontology源码](https://github.com/ontio/ontology)
5. [ontology-go-sdk源码](https://github.com/ontio/ontology-go-sdk)