CreditDB: A decentralized personal credit database system
===

A cybernet course project.

Author:
Yiwei Li (liyw19@mails.tsinghua.edu.cn)

Use
---

This project is highly based on the popular hyperledger tutorial (fabric-sample repo), you are highly recommended to read the tutorial and configure dependencies (docker, docker-compose, etc.) first: https://hyperledger-fabric.readthedocs.io/zh_CN/release-2.2/test_network.html

`source env.sh` - Specify Golang and peer path

`cd fabcar && ./startFabric.sh` - Enable the blockchain network, and the fabcar chaincode

`cd gin && go run gin.go` - Start the backend server, working at `0.0.0.0:8085`

Now you can access the accounts by sending GET or POST requests, e.g., using `curl`

`curl http://localhost:8085/hello` - Display project description

`curl http://localhost:8085/listAccount` - Show all accounts in Json array format

`curl http://localhost:8085/createAccount -X POST -d "name=CAR7811&field4=YiweiLi"` - Create a new account with specified fields. Return `Err=""` if successful

`curl http://localhost:8085/queryAccount -X POST -d "name=CAR7811"` - Query the fields of one specified account

Teardown
---

After finishing queries, close the server.

After finishing the blockchain testnet, close the network by `./networkDone.sh`