# TrustEx
## 1.项目架构
### 1.1 项目结构
```
(FinanceApp)TrustEx
├── README.md
├── asset-transfer-events(backend)
├── frontend

```
### 1.2 技术栈
- 后端：Node.js(v20.17.0) + Gin(go1.23.0)
- 前端：Vue.js + ElementUI
- 区块链：Hyperledger Fabric(latest)

## 2.项目启动
### 2.1 后端
部署完区块链网络之后，进入`asset-transfer-events/application-gateway-go`目录,执行以下命令：
```shell
go run main.go
```
### 2.2 前端
进入`frontend/finance_app/src`目录，执行以下命令：
```shell
npm install
npm run serve
```

## 3.项目功能
基于区块链的履约金融系统，实现了以下功能：
- 贷款/保险合同的创建、理赔、还款
- 资产的转让、查询

## 4.开发者
- [x] [陈其鹏](https://github.com/LeChatelier-Lenz)
- [x] [鲍熙元](https://github.com/LemonKozzet)
