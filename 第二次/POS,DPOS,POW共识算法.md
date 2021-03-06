> 笔记由 @无名 整理
> Date: 2018/06/28

---
### POW：Proof of Work，工作证明

- 符合要求的Hash由N个前导零构成
- 零的个数取决于网络的难度值
- 计算时间取决于机器的哈希运算速度
- 判断方法
  - 第一种，判断前N个是否为0
  - 第二种，判断hash是否小于等于target值
    - target为符合要求的最大hash值
- 查找hash的过程是一个死循环，其它数据不变的情况下，通过```nonce+1``` 改变hash，直到找到符合要求的hash
- 适合条件的hash，不止一个
- 寻找符合要求的hash是一个概率事件
- 节点占全网n%的算力，该节点有n%的概率找到适合要求的Hash
- 第一代共识机制,比特币的基础
- “按劳取酬”，你付出多少工作量，就会获得多少报酬
- 目前比特币、以太坊使用此算法

### POS：Proof of Stake，股权证明
- 点点币（PPC）的创新
- 类似股票/银行存款
- 想要获得更多的收益必须在线
- 节能，不用挖矿，不需要大量耗费电力和能源
- 持有币的人, 就有对应的权利, 持有的越多, 权利越大(收益就更多)

### DPOS：Delegated Proof of Stake，委任权益证明
- POS升级版
- 比特股（BTS）最先引入的
- 类似于董事会。董事会成员数量有限,由大家选举产生。被选中的董事会成员可以行使权利。
- 能耗更低
  - DPoS机制将节点数量进一步减少到101个
  - 会不会容易受DDOS攻击？？
- EOS使用此算法


---
### 个人观点：
- POW
  - 成本太高，随着算力的提升，消耗的电力更多
  - 矿池的出现，有点类似DPOS中见证人，也会趋同中心化
- POS
  - 如果初期参与人少，收益会过于集中，会趋同中心化
  - 后期随着币数量的集中度，也会趋同中心化
- DPOS
  - 见证人机制
    - 见证人数量有限
    - 见证人互相竞争来获得记账的工作
    - 见证人主动降低他们获得的收入，获得投票
    - 见证人工资支付在更多工作而获得选票，例如进行比特股的市场推广、法务
    - 更像一个24小时不间断的股东大会，股东们可以在任意时间通过投票改变公司的组织架构
  - 良好的生态，利于发展
  - 利益的关系，会使见证人在经过一段时间后趋于固定，趋同中心化



---
### 参考：
- [共识算法（POW,POS,DPOS,PBFT）介绍和心得](https://blog.csdn.net/lsttoy/article/details/61624287)
- [DPoS（委托权益证明机制）官方共识机制详解——BTS、EOS](https://blog.csdn.net/lsttoy/article/details/80041033)
- [POW , POS 与 DPOS 一切都为了共识](https://www.jianshu.com/p/f99e8fe57c9a  )
- [解读《EOS.IO技术白皮书》](https://www.jianshu.com/p/bc489db794ce)
- [DPOS委托权益证明 vs POW工作量证明](https://zhuanlan.zhihu.com/p/28127511)
- [授权股权证明机制 (DPOS)白皮书翻译修订版](http://vdisk.weibo.com/s/yUM-q4I09jbOT)