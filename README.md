CreditDB: A decentralized personal credit database system
===

A cybernet course project.

Author:
Yiwei Li (liyw19@mails.tsinghua.edu.cn)


Dependencies configuration
---

#### Docker

This project is highly based on the popular hyperledger tutorial (fabric-sample repo), you are highly recommended to read the tutorial and configure dependencies (docker, docker-compose, etc.): https://hyperledger-fabric.readthedocs.io/zh_CN/release-2.2/prereqs.html

Specifically, install docker and docker-compose by apt: `sudo apt install docker docker-compose`

Add current user to docker group: `sudo gpasswd -a ${USER} docker`

Now you are able to run `docker ps` as a non-priviledged user. If you are running within a tmux session, make sure you close all previous tmux windows and restart a new window. If you still cannot continue, restart this server.

You are suggested to configure the docker source (https://www.jianshu.com/p/405fe33b9032)

#### Golang

This project uses Go for chaincode (Hyperledger Go SDK) and back-end (Gin framework). Install them following: https://golang.org/doc/install

Modify the `env.sh` for the go binary path.

Use
---

Make sure you have finished configuring dependencies.

`source env.sh` - Specify Golang and peer path

`cd creditdb && ./startFabric.sh` - Enable the blockchain network, and the creditdb chaincode

`cd gin && go run gin.go` - Start the backend server, working at `0.0.0.0:8085`

Now you can access the accounts by sending GET or POST requests, e.g., using `curl`

`curl http://localhost:8085/hello` - Display project description

`curl http://localhost:8085/listTx` - Show all accounts in Json array format

`curl http://localhost:8085/createTx -X POST -d "id=123&value=1.00&SenderName=Bank&RecverName=Alice"` - Create a new account with specified fields. Return `Err=""` if successful

`curl http://localhost:8085/queryUser -X POST -d "name=Alice"` - Query the fields of one specified account

Teardown
---

After finishing queries, close the server.

After finishing the blockchain testnet, close the network by `./networkDown.sh`