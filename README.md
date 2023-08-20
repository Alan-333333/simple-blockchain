# Simple Blockchain

这是一个用于学习目的的简单区块链实现。

## 简介

该项目使用 Go 语言实现了区块链的基本功能,包含:

- 区块和交易数据结构
- 地址和钱包管理 
- 挖矿和工作量证明
- 链式存储区块
- 简单的网络通信
- 命令行界面

## 概述

该区块链具有以下关键功能：

- 实现基本的工作量证明算法
- 支持创建交易
- 维护交易池
- 验证和挖掘块
- 实现简单的P2P网络用于块和交易传播
- 提供与命令行交互的CLI

## 运行

1. 安装Go 1.19或更高版本

2. 下载代码

```
git clone https://github.com/Alan-333333/simple-blockchain.git
```

3. 进入项目目录

```
cd simple-blockchain
```

4. 安装依赖

```
go mod tidy
```

5. 运行节点

```
go run main.go
```

## 用法

支持以下命令：

- `printBlockchain` - 打印区块链中的所有块
- `getBlock <hash>` - 打印块
- `createGenesisBlock` - 创建创世块
- `createWallet` - 创建一个新的钱包
- `getBalance <address>` - 获取钱包地址的余额
- `sendTransaction <from> <to> <amount>` - 创建并发送交易
- `connectNode <ip> <port>` - 连接到节点
- `addWalletBalance <address> <amount>` - 向钱包添加余额

## 快速开始

1. 运行区块链节点：
<!---->

```
go run main.go
```
2.  运行另一个节点并连接到第一个节点：

<!---->

```
go run main.go -port 3001
connectNode 127.0.0.1 3000
```
3.  创建钱包：

<!---->

```
createWallet
createWallet
```
4.  添加余额：

<!---->

```
addWalletBalance 1K4nFZNxmHRRwfM4E9S8SXPQcTcayxaeKj 50
addWalletBalance 1JLfCguhUBui6MWQ4vNFgEktn87E9V6F8Q 100
```
5. 验证余额：

<!---->

```
getBalance 1K4nFZNxmHRRwfM4E9S8SXPQcTcayxaeKj
getBalance 1JLfCguhUBui6MWQ4vNFgEktn87E9V6F8Q
```
6.  发送交易：

<!---->

```
sendTransaction 1K4nFZNxmHRRwfM4E9S8SXPQcTcayxaeKj 1JLfCguhUBui6MWQ4vNFgEktn87E9V6F8Q 50
```

## 贡献

欢迎提 issue 和 PR 为项目作出贡献!

## 版权

该项目采用 MIT 许可证,详情请参阅 LICENSE 文件。

