
> 方便测试，钱包地址不清除

#### 编译程序
```
# 删除旧的区块链库
rm -rf blockchain.db

go build -o bc main.go

```
#### 打印钱包地址

```
./bc listaddresses
1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72
1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf
1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ
```

#### 创建创世区块
```
./bc createblockchain -address 1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72
```

#### 查看余额
```
./bc getbalanceall

地址:1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf的余额为：0
地址:1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ的余额为：0
地址:1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72的余额为：10

```

#### 测试单笔转帐(挖矿奖励10)
```
./bc send -from '["1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72"]' -to '["1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ"]' -amount '["4"]'
```
- 1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72 ： 10 - 4 + 10 = 16

- 1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ ： 4

> 测试余额

```
地址:1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf的余额为：0
地址:1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ的余额为：4
地址:1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72的余额为：16

```
#### 测试多笔转账
```
./bc send -from '["1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72","1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72"]' -to '["1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ","1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf"]' -amount '["5","5"]'
```

- 1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf   0 + 5 = 5
- 1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ   4 + 5 = 9
- 1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72   16 - 5 - 5 + 10 = 16

> 测试余额

```
地址:1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf的余额为：5
地址:1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ的余额为：9
地址:1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72的余额为：16
```
#### 继续测试多笔转账

```
./bc send -from '["1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72","1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72"]' -to '["1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ","1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf"]' -amount '["3","4"]'
```
> 测试余额

```
./bc getbalanceall

地址:1PE646PctnSH5hRfUWMQGqJerm8Emx9gCf的余额为：9
地址:1Hp2wo6jiGghei1eDMX8Y5v44JdNXxaaKZ的余额为：12
地址:1MWd2iQoYyjQFcUW1GU48qxCNqoici7p72的余额为：19

```

