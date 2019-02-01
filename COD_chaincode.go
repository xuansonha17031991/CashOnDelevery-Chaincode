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
	Balance    int    `json:"balance"`
}

type Customer struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Number     string `json:"number"`
	Email      string `json:"email"`
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
	case "createCustomer":
		return t.createCustomer(stub, args)
	case "createSeller":
		return t.createSeller(stub, args)
	case "query":
		return t.query(stub, args)
	case "createBalance":
		return t.createBalance(stub, args)
	case "transferMoney":
		return t.transferMoney(stub, args)
	case "createOrder":
		return t.createOrder(stub, args)
	default:
		fmt.Println("Invoke did not find function: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

//create customer information
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

	name := args[0]
	location := args[1]
	number := args[2]
	email := args[3]

	//convert variable to json
	objectType := "Customer"
	customer := &Customer{objectType, name, location, number, email}
	customer_to_byte, err := json.Marshal(customer)
	if err != nil {
		return shim.Error(err.Error())
	}

	//save to database
	err = stub.PutPrivateData("CODcollection", name, customer_to_byte)
	if err != nil {
		return shim.Error(err.Error())
	}

	//create and save key
	indexName := "name~number"
	customerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{customer.Name, customer.Number})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutPrivateData("CODcollection", customerNameIndexKey, value)

	return shim.Success(nil)
}

//create seller information
func (t *COD_chaincode) createSeller(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	name := args[0]
	asset := args[1]
	quantity, q_err := strconv.Atoi(args[2])
	price, p_err := strconv.Atoi(args[3])

	if q_err != nil {
		return shim.Error("quantity must be a number")
	}
	if p_err != nil {
		return shim.Error("price must be a number")
	}

	//convert variable to json
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

//query data
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

//create delivery information
func (t *COD_chaincode) createBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("expecting 2 argument, name and balance")
	}

	name := args[0]
	balance, err_owner_balance := strconv.Atoi(args[1])
	if err_owner_balance != nil {
		return shim.Error("balance must be a number")
	}

	//convert to json
	objectType := "Balance"
	owner := &Balance{objectType, name, balance}
	owner_to_byte, err := json.Marshal(owner)
	if err != nil {
		return shim.Error(err.Error())
	}

	//save to ledger
	err = stub.PutPrivateData("CODcollection", name, owner_to_byte)
	if err != nil {
		return shim.Error(err.Error())
	}

	//create and save key
	indexName := "name~balance"
	balanceNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{owner.Name, strconv.Itoa(owner.Balance)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutPrivateData("CODcollection", balanceNameIndexKey, value)

	return shim.Success(nil)
}

//transfer money to new owner
func (t *COD_chaincode) transferMoney(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("expecting 3 parameter precent owner, money, new owner")
	}

	old_owner := Balance{}
	new_owner := Balance{}
	old_owner_name := args[0]
	mortgage, err_mortgage := strconv.Atoi(args[1])
	if err_mortgage != nil {
		return shim.Error("balance isn't a number")
	}
	new_owner_name := args[2]

	//get old owner's information
	old_owner_as_byte, err := stub.GetPrivateData("CODcollection", old_owner_name)
	if err != nil {
		return shim.Error("cannot get owner's infor")
	} else if old_owner_as_byte == nil {
		return shim.Error("owner doesn't exist")
	}

	//unmarshal to old_owner variable
	err = json.Unmarshal(old_owner_as_byte, &old_owner)
	if err != nil {
		return shim.Error("loi khong the unmarshal")
	}

	//check old_owner balance
	if old_owner.Balance < mortgage {
		return shim.Error("present owner does not enough balance")
	}
	old_owner.Balance = old_owner.Balance - mortgage

	new_info_old_owner_as_byte, err := json.Marshal(old_owner)

	err = stub.PutPrivateData("CODcollection", old_owner_name, new_info_old_owner_as_byte)
	if err != nil {
		return shim.Error("cannot save new info of old owner")
	}

	//get info of virtual account
	new_owner_as_byte, err := stub.GetPrivateData("CODcollection", new_owner_name)
	if err != nil {
		return shim.Error("cannot get info of new owner")
	}
	err = json.Unmarshal(new_owner_as_byte, &new_owner)
	if err != nil {
		return shim.Error("cannot unmarshal new owner")
	}
	new_owner.Balance = new_owner.Balance + mortgage
	new_owner_new_info_as_byte, err := json.Marshal(new_owner)
	if err != nil {
		return shim.Error("cannot marshal new info of new owner")
	}
	err = stub.PutPrivateData("CODcollection", new_owner_name, new_owner_new_info_as_byte)
	if err != nil {
		return shim.Error("cannot put new data of new owner")
	}

	fmt.Println("transfer successful")
	return shim.Success(nil)
}

//create order information
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

func (t *COD_chaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var owner_new_info Balance
	var jsRespon string

	if len(args) != 1 {
		return shim.Error("expecting name of user")
	}

	name := args[0]

	//get thong tin cua user
	old_owner_as_byte, err := stub.GetPrivateData("CODcollection", name)
	if err != nil {
		return shim.Error(err.Error())
	}

	//transfer data into new variable
	err = json.Unmarshal([]byte(old_owner_as_byte), &owner_new_info)
	if err != nil {
		jsRespon = "{\"Error\":\"Failed to decode JSON of: " + name + "\"}"
		return shim.Error(jsRespon)
	}

	//xoa du lieu
	err = stub.DelPrivateData("CODcollection", name)
	if err != nil {
		return shim.Error("data does not exist")
	}

	//tao lai gia tri key
	indexName := "name~balance"
	balanceNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{owner_new_info.Name, strconv.Itoa(owner_new_info.Balance)})
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.DelPrivateData("CODcollection", balanceNameIndexKey)
	if err != nil {
		return shim.Error("cannot delete key")
	}

	return shim.Success(nil)
}
