package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type COD_chaincode struct {
}

type Seller struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Asset      string `json:"asset"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
}

type Balance struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Balance    string `json:"balance"`
}

type Customer struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Number     string `json:"number"`
	Email      string `json:"email"`
}

type Coupon struct {
	ObjectType string `json:"docType"`
	CouponID   string `json:"couponid"`
	Customer   string `json:"customer"`
	Seller     string `json:"seller"`
	Asset      string `json:"asset"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
}

type Order struct {
	ObjectType string `json:"docType"`
	OrderID    string `json:"orderid"`
	Customer   string `json:"customer"`
	Seller     string `json:"seller"`
	Delivery   string `json:"delivery"`
	Asset      string `json:"asset"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
	Status     string `json:"status"`
}

/*main*/
func main() {
	err := shim.Start(new(COD_chaincode))
	if err != nil {
		fmt.Printf("cannot initiate COD chaincode: %s", err)
	}
}

// init chaincode
func (t *COD_chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke
func (t *COD_chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke is running" + function)

	switch function {
	case "createSeller":
		return t.createSeller(stub, args)
	case "querySeller":
		return t.query(stub, args)
	case "createBalance":
		return t.createBalance(stub, args)
	case "createCustomer":
		return t.createCustomer(stub, args)
	default:
		fmt.Println("Invoke did not find function: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

//khoi tao customer
func (t *COD_chaincode) createCustomer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("expecting 4 argument")
	}

	if len(args[0]) == 0 {
		return shim.Error("Customer's name must be declare")
	}
	if len(args[1]) == 0 {
		return shim.Error("Customer's location must be declare")
	}
	if len(args[2]) == 0 {
		return shim.Error("Customer's number must be declare")
	}
	if len(args[3]) == 0 {
		return shim.Error("Customer's email must be declare")
	}

	//gan bien
	name := args[0]
	location := args[1]
	number := args[2]
	email := args[3]

	//convert qua son
	objectType := "Customer"
	customer := &Customer{objectType, name, location, number, email}
	customer_to_byte, err := json.Marshal(customer)
	if err != nil {
		return shim.Error(err.Error())
	}

	//save to database
	err = stub.PutPrivateData("CODcolection", name, customer_to_byte)
	if err != nil {
		return shim.Error(err.Error())
	}

	//create key
	indexName := "name~number"
	customerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{customer.Name, customer.Number})
	if err != nil {
		return shim.Error(err.Error())
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("CODcollection", customerNameIndexKey, value)

	return shim.Success(nil)
}

//khoi tao thong tin seller
func (t *COD_chaincode) createSeller(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//args: name of seller, name of asset, count, price
	if len(args) != 4 {
		return shim.Error("there must be 4 argument")
	}

	if len(args[0]) == 0 {
		return shim.Error("name of seller must be declare")
	}
	if len(args[1]) == 0 {
		return shim.Error("name of asset must be declare")
	}
	if len(args[2]) == 0 {
		return shim.Error("quantity must be declare")
	}
	if len(args[0]) == 0 {
		return shim.Error("price must be declare")
	}

	//gan bien
	name := args[0]
	asset := args[1]
	quantity, q_err := strconv.Atoi(args[2])
	price, p_err := strconv.Atoi(args[3])

	//kiem tra du lieu
	if q_err != nil {
		return shim.Error("quantity must be a number")
	}
	if p_err != nil {
		return shim.Error("price must be a number")
	}

	//convert into json
	objectType := "Seller"
	seller := &Seller{objectType, name, asset, quantity, price}
	seller_to_byte, err := json.Marshal(seller)
	if err != nil {
		return shim.Error(err.Error())
	}

	//save to database
	err = stub.PutPrivateData("CODcollection", name, seller_to_byte)
	if err != nil {
		return shim.Error(err.Error())
	}

	//create index key
	indexName := "name~asset"
	assetNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{seller.Name, seller.Asset})
	if err != nil {
		return shim.Error(err.Error())
	}

	//save index
	value := []byte{0x00}
	stub.PutPrivateData("CODcollection", assetNameIndexKey, value)

	return shim.Success(nil)
}

//Get database of seller
func (t *COD_chaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("sai roi, nhap ten")
	}

	name = args[0]
	valAsBytes, err := stub.GetPrivateData("CODcollection", name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsBytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsBytes)
}

//Create the information of delivery
func (t *COD_chaincode) createBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("expecting 2 argument, name and balance")
	}

	//Create information
	name := args[0]
	balance := args[1]

	//convert into json
	objectType := "Balance"
	owner := &Balance{objectType, name, balance}
	owner_to_byte, err := json.Marshal(owner)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Store data
	err = stub.PutPrivateData("CODcollection", name, owner_to_byte)
	if err != nil {
		return shim.Error(err.Error())
	}

	//create key
	indexName := "name~balance"
	balanceNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{owner.Name, owner.Balance})
	if err != nil {
		return shim.Error(err.Error())
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("CODcollection", balanceNameIndexKey, value)

	return shim.Success(nil)
}

//
func (t *COD_chaincode) transferMoney(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("expecting 3 parameter precent owner, money, new owner")
	}

	owner := args[0]
	money := args[1]
	new_owner := args[3]
	fmt.Println("transfering money")

	precent_owner_as_byte, err := stub.GetPrivateData("CODcollection", owner)
	if err != nil {
		return shim.Error("cannot get owner's infor")
	} else if precent_owner_as_byte == nil {
		return shim.Error("owner doesn't exist")
	}

	transfer := Balance{}
	err = json.Unmarshal(precent_owner_as_byte, &transfer)
	if err != nil {
		return shim.Error(err.Error())
	}

	//check balance
	if transfer.Balance < money {
		return shim.Error("the owner's balance is not enought")
	}

	transfer.Name = new_owner

	new_ownerAsByte, _ := json.Marshal(transfer)
	err = stub.PutPrivateData("CODcollection", new_owner, new_ownerAsByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("transfer successful")
	return shim.Success(nil)
}

//create coupon
func (t *COD_chaincode) createCoupon(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6 {
		return shim.Error("there is not enought argument, expecting 5")
	}

	coupon_id := args[0]
	customer := args[1]
	seller := args[2]
	asset := args[3]
	quantity, err := strconv.Atoi(args[4])
	price, err := strconv.Atoi(args[5])

	objectType := "Coupon"
	coupon := &Coupon{objectType, coupon_id, customer, seller, asset, quantity, price}
	coupon_as_byte, err := json.Marshal(coupon)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData("CODcollection", coupon_id, coupon_as_byte)
	if err != nil {
		return shim.Error(err.Error())
	}

	//create key
	indexName := "coupon_id~customer"
	couponNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{coupon.CouponID, coupon.Customer})
	if err != nil {
		return shim.Error(err.Error())
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("CODcollection", couponNameIndexKey, value)

	return shim.Success(nil)
}

//create order
func (t *COD_chaincode) createOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 8 {
		return shim.Error("expecting 8 argument")
	}

	id := args[0]
	customer := args[1]
	seller := args[2]
	delivery := args[3]
	asset := args[4]
	quantity, err_qu := strconv.Atoi(args[5])
	price, err_pr := strconv.Atoi(args[6])
	status := args[7]

	if err_qu != nil {
		return shim.Error("quantity must be a number")
	}
	if err_pr != nil {
		return shim.Error("price must be a number")
	}

	objectType := "Order"
	order := &Order{objectType, id, customer, seller, delivery, asset, quantity, price, status}
	order_to_byte, err_or := json.Marshal(order)
	if err_or != nil {
		return shim.Error(err_or.Error())
	}

	err_or = stub.PutPrivateData("CODcollection", id, order_to_byte)
	if err_or != nil {
		return shim.Error(err_or.Error())
	}

	//create key
	indexName := "id~name"
	orderNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{order.OrderID, order.Customer})
	if err != nil {
		return shim.Error(err.Error())
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("CODcollection", orderNameIndexKey, value)

	return shim.Success(nil)
}
