#Cash On Delivery Chain-code
##Description
this is smart contract if the Cash On Delivery (COD) system 
##Prerequisites
you need to install fabric-examples version 1.4 on you machine, but we re-comment it should be installed into go path directory. If this is a first time run (curl -sSL http://bit.ly/2ysbOFE | bash -s 1.4.0-rc2) if not, just clone fabric-examples by run (git clone https://github.com/hyperledger/fabric-samples.git) command
##Installation
- Go to fabric-samples's chain-code directory  
- Open your terminal and run this command to install
```bash
git clone 
```
- Now you're ready
##Usage
- Go to first-network directory of fabric-sample
- Open terminal and run ./byfn.sh down to to clean up any previous environments:
```
./byfn.sh down
```
- Startup your network:
```
./byfn.sh up -c [your channel's name]
```
###Install and instantiate chain-code with a collection
- BYFN network is started, enter the CLI container.
```
docker exec -it cli bash
```
- the result will be 
```
root@81eac8493633:/opt/gopath/src/github.com/hyperledger/fabric/peer#
```
- Run this command
```
peer chaincode install -n COD -v 1.0 -p github.com/chaincode/COD/COD_chaincode/
```
- If success the result will be:
```
install -> INFO 003 Installed remotely response:<status:200 payload:"OK" >
```
###Instantiate the chaincode
- We use the peer chaincode instantiate command to instantiate the COD chaincode and collection.json file for configuration
```
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```
```
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n COD -v 1.0 -c '{"Args":["init"]}' -P "OR('Org1MSP.member','Org2MSP.member')" --collections-config  $GOPATH/src/github.com/chaincode/COD/COD_chaincode/collection.json
```
- If the instantiation completes successfully you should see:
```
[chaincodeCmd] checkChaincodeCmdParams -> INFO 001 Using default escc
[chaincodeCmd] checkChaincodeCmdParams -> INFO 002 Using default vscc
```
- Now you can use peer command to invoke functions on this smartcontract. Have a nice day!!!
###Store customer's data
##Contributing
You're welcome to make a pull requests
##Support
- Ha Xuan Son - hxson@ctuet.edu.vn
- Nguyen Quoc Nghiep - nqnghiep.ktpm0115@student.ctuet.edu.vn
- Le Hai Trieu - lhtrieu.ktpm0115@student.ctuet.edu.vn
##Roadmap
We will develop a COD system that applies blockchain technology based on this chaincode
##Authors and acknowledment
Ha Xuan Son  - Team leader
Nguyen Quoc Nghiep - Develop
Le Hai Trieu - Develop
##Status
This project in processing