```
./bc

Usage:
  createwallet - 创建钱包
  listaddresses - 打印钱包地址
  createblockchain -address ADDRESS - 创建区块链
  getbalance -address ADDRESS - 获取地址的余额
  getbalanceall - 打印所有钱包地址的余额
  printchain - 打印区块链中的所有区块数据
  send -from FROM -to TO -amount AMOUNT 转账
  reindexutxo - 重建UTXO set
  printutxo - 打印所有的UTXO set

```
#### 2018/07/10 第六次作业
#### 2018/07/17 第七次作业
##### 截止目前实现功能如下：
- 创建钱包
  - 公钥、私钥
  - 椭圆曲线加密
  - Base58Encode(version + ripemd160(sha256(PubKey)) + sha256(sha256(PubKeyHash))[:addressChecksumLen])
  - 钱包地址可以到 http://blockchain.info/ 验证
- 查看所有钱包地址
- 创建创世区块
  - 交易输入(Vin/TXInput)
    - 只有一条
    - txID=0
    - vout=-1
  - 交易输出(Vout/TXOutput)
    - 通常只有一条
    - 区块链初始金额
- 多笔交易
  - 支持点对点交易
    - ./bc send -from '["1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72"]' -to '["1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ"]' -amount '["4"]'
  - 支持多对多交易(多笔交易在一个区块中打包)
    - ./bc send -from '["1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72","1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72"]' -to '["1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ","1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf"]' -amount '["5","5"]'
  - 挖矿
    - 工作量证明(POW)
      - Merkle 树
    - 奖励(目前实现比较简单)
  - UTXO
  - UTXO集
    - 将所有未花费的交易输出存储到数据库中，提高性能
  - 交易数字签名及验证
- 查看指定余额
- 查看所有钱包地址余额
- 查看所有区块链中数据
- 重建UTXO集
  - 重建会比较花费时间
- 查看所有UTXO集中的内容
- 多节点通讯功能
  - 主节点
  - 钱包节点
  - 矿工节点

数据库：[BoltDB](https://github.com/boltdb/bolt)
[ripemd160](https://github.com/golang/crypto)

- [第六次测试文档](https://github.com/kongyixueyuan/renweiqing/blob/master/%E7%AC%AC%E5%85%AD%E6%AC%A1/test.md)
- [第七次测试文档](https://github.com/kongyixueyuan/renweiqing/blob/master/%E7%AC%AC%E4%B8%83%E6%AC%A1/test.md)
- [第八次作业说明(测试说明参考第七次测试文档)](https://github.com/kongyixueyuan/renweiqing/blob/master/%E7%AC%AC%E4%B8%83%E6%AC%A1/test.md)

----
#### 钱包地址生成方式

![image](http://upload-images.jianshu.io/upload_images/127313-6aa6cff5d863d496.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

上图中Checksum生成方式不对
```
Version  Public key hash                           Checksum
00       62E907B15CBF27D5425399EBF6F0FB50EBB88F18  C29B7D93
```
> Base58Encode(version + ripemd160(sha256(PubKey)) + sha256(sha256(PubKeyHash))[:addressChecksumLen])

#### 交易签名

![image](http://upload-images.jianshu.io/upload_images/127313-ec45a7fca855f2e0.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


#### Merkle 树

![image](http://upload-images.jianshu.io/upload_images/127313-9c708d3c3d6a19c2.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)



