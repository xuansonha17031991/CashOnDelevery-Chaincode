package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type COD_chaincode struct {
}

type Asset struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Asset      string `json:"asset"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
}

type OrderHash struct {
	ObjectType string `json:"docType"`
	OrderID    string `json:"orderid"`
	AssetHash  string `json:"assethash"`
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
	AssetName  string `json:"assetname"`
	Detail     string `json:"detail"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
	Status     string `json:"status"`
}

type ImageAsByte struct {
	ObjectType  string `json:"docType"`
	OrderID     string `json:"orderid"`
	ImageAsByte byte   `json:"imageasbyte"`
}

type Delivery struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Price      int    `json:"price"`
	Distance   string `json:"distance"`
	Time       string `json:"time"`
}

type VerifyShipper struct {
	ObjectType string `json:"docType"`
	OrderID    string `json:"orderid`
	Hash       string `json:"hash"`
	Status     string `json:"status"`
	Location   string `json:"location"`
}

type LimitTime struct {
	ObjectType string `json:"docType"`
	OrderID    string `json:"orderid"`
	SellerID   string `json:"sellerid"`
	DeliveryID string `json:"deliveryid"`
	Time       string `json:"limittime"`
	Day        string `json:"day"`
}

/*main*/
func main() {
	err := shim.Start(new(COD_chaincode))
	if err != nil {
		fmt.Printf("cannot initiate COD chaincode: %s", err)
	}
}

// Init chaincode
func (t *COD_chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke
func (t *COD_chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke is running" + function)

	switch function {
	case "encrypAsset":
		return t.encrypAsset(stub, args)
	case "createAsset":
		return t.createAsset(stub, args)
	case "createAssetHash":
		return t.createAssetHash(stub, args)
	case "createBalance":
		return t.createBalance(stub, args)
	case "createCustomer":
		return t.createCustomer(stub, args)
	case "createDelivery":
		return t.createDelivery(stub, args)
	case "createOrder":
		return t.createOrder(stub, args)
	case "dealLimitTime":
		return t.dealLimitTime(stub, args)
	case "delete":
		return t.delete(stub, args)
	// case "imageToByte":
	// 	return t.imageToByte(stub, args)
	case "query":
		return t.query(stub, args)
	case "transferMoney":
		return t.transferMoney(stub, args)
	case "verifyShipper":
		return t.verifyShipper(stub, args)

	default:
		fmt.Println("Invoke did not find function: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

//create customer information
func (t *COD_chaincode) createCustomer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createCustomer function ===============")
	start := time.Now()
	time.Sleep(time.Second)
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
	err = stub.PutPrivateData("customerCollection", name, customer_to_byte)
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
	stub.PutPrivateData("customerCollection", customerNameIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("\nfunction createCustomer")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("execute time: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end createCustomer function ===============")
	return shim.Success(nil)
}

//create asset
func (t *COD_chaincode) createAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createAsset function ===============")
	start := time.Now()
	time.Sleep(time.Second)
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
	seller := &Asset{objectType, name, asset, quantity, price}
	seller_to_byte, err := json.Marshal(seller)
	if err != nil {
		return shim.Error(err.Error())
	}

	//save to database
	err = stub.PutPrivateData("assetCollection", name, seller_to_byte)
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
	stub.PutPrivateData("assetCollection", assetNameIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("\nfunction createAsset")
	fmt.Printf("time start: %s", start.String())
	fmt.Printf("time end: %s", end.String())
	fmt.Println("time execute: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end createAsset function ===============")
	return shim.Success(nil)
}

//query data
func (t *COD_chaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start query function ===============")
	start := time.Now()
	time.Sleep(time.Second)
	var name, jsonResp string
	var err error

	if len(args) != 2 {
		return shim.Error("index of object is invalid")
	}

	name = args[0]
	valAsBytes, err := stub.GetPrivateData(args[1], name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsBytes == nil {
		jsonResp = "{\"Error\":\"object does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("function query")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end query function ===============")

	return shim.Success(valAsBytes)
}

//create delivery information
func (t *COD_chaincode) createBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createBalance function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	if len(args) != 3 {
		return shim.Error("expecting 3 argument, name and balance")
		// return "expecting 3 argument, name and balance"
	}

	name := args[0]
	balance, err_owner_balance := strconv.Atoi(args[1])
	if err_owner_balance != nil {
		return shim.Error("balance must be a number")
		// return "Error"
	}
	collection := ""
	switch args[2] {
	case "Org1":
		collection = "balanceOrg1Collection"
	case "Org2":
		collection = "balanceOrg2Collection"
	case "mortgage":
		collection = "mortgageCollection"
	}

	//convert to json
	objectType := "Balance"
	owner := &Balance{objectType, name, balance}
	owner_to_byte, err := json.Marshal(owner)
	if err != nil {
		return shim.Error(err.Error())
		// return "Error"
	}

	//save to ledger
	err = stub.PutPrivateData(collection, name, owner_to_byte)
	if err != nil {
		return shim.Error(err.Error())
		// return "Error"
	}

	//create and save key
	indexName := "name~balance"
	balanceNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{owner.Name, strconv.Itoa(owner.Balance)})
	if err != nil {
		return shim.Error(err.Error())
		// return " Error"
	}
	value := []byte{0x00}
	stub.PutPrivateData(collection, balanceNameIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction createBalance")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end createBalance function ===============")

	return shim.Success(nil)
}

func (t *COD_chaincode) createDelivery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createDelivery function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	//check length of data
	if len(args) != 5 {
		return shim.Error("expecting 5 argument")
	}

	//definite data variable
	name := args[0]
	location := args[1]
	price, errPrice := strconv.Atoi(args[2])
	if errPrice != nil {
		return shim.Error("prive must be a number")
	}
	distance := args[3]
	Dtime := args[4]
	ObjectType := "Delivery"

	delivery := &Delivery{ObjectType, name, location, price, distance, Dtime}

	//marshal delivery to byte
	deliveryAsByte, errDelivery := json.Marshal(delivery)
	if errDelivery != nil {
		return shim.Error("cannot marshal delivery's data")
	}

	//put data to ledger
	errDelivery = stub.PutPrivateData("deliveryCollection", name, deliveryAsByte)
	if errDelivery != nil {
		return shim.Error("cannot put private data of delivery")
	}

	//create index key
	indexKey := "name"
	deliveryIndexKey, errDeliveryIndexKey := stub.CreateCompositeKey(indexKey, []string{delivery.Name, delivery.Location, strconv.Itoa(delivery.Price), delivery.Distance, delivery.Time})
	if errDeliveryIndexKey != nil {
		return shim.Error("cannot create index key of delivery")
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("deliveryCollection", deliveryIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction createBalance")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end createDelivery function ===============")
	time.Sleep(time.Second)

	return shim.Success(nil)
}

//transfer money to new owner
func (t *COD_chaincode) transferMoney(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start transferMoney function ===============")
	start := time.Now()
	time.Sleep(time.Second)
	if len(args) != 5 {
		return shim.Error("expecting 5 arguments")
	}

	old_owner := Balance{}
	new_owner := Balance{}
	old_owner_name := args[0]
	mortgage, err_mortgage := strconv.Atoi(args[1])
	if err_mortgage != nil {
		return shim.Error("balance isn't a number")
	}
	new_owner_name := args[2]
	collection1 := args[3]
	collection2 := args[4]

	//get old owner's information
	old_owner_as_byte, err := stub.GetPrivateData(collection1, old_owner_name)
	if err != nil {
		return shim.Error("cannot get owner's infor")
	} else if old_owner_as_byte == nil {
		return shim.Error("owner doesn't exist")
	}

	//unmarshal to old_owner variable
	err = json.Unmarshal(old_owner_as_byte, &old_owner)
	if err != nil {
		return shim.Error("cannot unmarshal data")
	}

	//check old_owner balance
	if old_owner.Balance < mortgage {
		return shim.Error("present owner does not enough balance")
	}
	old_owner.Balance = old_owner.Balance - mortgage

	new_info_old_owner_as_byte, err := json.Marshal(old_owner)

	err = stub.PutPrivateData(collection2, old_owner_name, new_info_old_owner_as_byte)
	if err != nil {
		return shim.Error("cannot save new info of old owner")
	}

	//get info of virtual account
	new_owner_as_byte, err := stub.GetPrivateData(collection2, new_owner_name)
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
	err = stub.PutPrivateData(collection2, new_owner_name, new_owner_new_info_as_byte)
	if err != nil {
		return shim.Error("cannot put new data of new owner")
	}

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("\nfunction transferMoney")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end transferMoney function ===============")
	return shim.Success(nil)
}

//create order information
func (t *COD_chaincode) createOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createOrder function ===============")
	start := time.Now()
	time.Sleep(time.Second)
	if len(args) != 9 {
		return shim.Error("expecting 9 argument")
	}

	id := args[0]
	customer := args[1]
	seller := args[2]
	delivery := args[3]
	assetname := args[4]
	detail := args[5]
	quantity, err_qu := strconv.Atoi(args[6])
	price, err_pr := strconv.Atoi(args[7])
	status := args[8]

	if err_qu != nil {
		return shim.Error("quantity must be a number")
	}
	if err_pr != nil {
		return shim.Error("price must be a number")
	}

	objectType := "Order"

	order := &Order{objectType, id, customer, seller, delivery, assetname, detail, quantity, price, status}
	order_to_byte, err_or := json.Marshal(order)
	if err_or != nil {
		return shim.Error(err_or.Error())
	}

	err_or = stub.PutPrivateData("orderCollection", id, order_to_byte)
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
	stub.PutPrivateData("orderCollection", orderNameIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction createOrder")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end createOrder function ===============")

	return shim.Success(nil)
}

func (t *COD_chaincode) createAssetHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createAssetHash function ===============")
	start := time.Now()
	time.Sleep(time.Second)
	if len(args) != 6 {
		return shim.Error("expting 6 parameters")
	}
	OrderID := args[0]
	sellerId := args[1]
	asset := args[2]
	detail := args[3]
	quantity, err1 := strconv.Atoi(args[4])
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	price, err2 := strconv.Atoi(args[5])
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	ObjectType := "AssetHash"
	hash := sha256.New()
	hash.Write([]byte(sellerId + asset + detail + string(quantity) + string(price)))
	md := hash.Sum(nil)
	asset_hash := hex.EncodeToString(md)

	AssetHash := &OrderHash{ObjectType, OrderID, asset_hash}
	AssetHashToByte, err := json.Marshal(AssetHash)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData("assetHashCollection", OrderID, AssetHashToByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	//create key
	indexName := "OrderID~Hash"
	orderHashIndexKey, errKey := stub.CreateCompositeKey(indexName, []string{ObjectType, OrderID, asset_hash})
	if errKey != nil {
		return shim.Error(errKey.Error())
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("assetHashCollection", orderHashIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction createAssetHash")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end createAssetHash function ===============")
	return shim.Success(nil)
}

func (t *COD_chaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var owner_new_info Balance
	var jsRespon string

	if len(args) != 2 {
		return shim.Error("expecting name of user")
	}

	name := args[0]

	//get user's information
	old_owner_as_byte, err := stub.GetPrivateData(args[1], name)
	if err != nil {
		return shim.Error(err.Error())
	}

	//transfer data into new variable
	err = json.Unmarshal([]byte(old_owner_as_byte), &owner_new_info)
	if err != nil {
		jsRespon = "{\"Error\":\"Failed to decode JSON of: " + name + "\"}"
		return shim.Error(jsRespon)
	}

	//delete data
	err = stub.DelPrivateData(args[1], name)
	if err != nil {
		return shim.Error("data does not exist")
	}

	//create key
	indexName := "name~balance"
	balanceNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{owner_new_info.Name, strconv.Itoa(owner_new_info.Balance)})
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.DelPrivateData(args[1], balanceNameIndexKey)
	if err != nil {
		return shim.Error("cannot delete key")
	}

	return shim.Success(nil)
}

func (t *COD_chaincode) verifyShipper(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start verifyShipper function ===============")
	start := time.Now()
	time.Sleep(time.Second)
	if len(args) != 3 {
		return shim.Error("there must be 3 arguments")
	}

	orderHash := OrderHash{}

	id := args[0]
	//get value as byte
	valAsBytes, err := stub.GetPrivateData("assetHashCollection", id)
	if err != nil {
		return shim.Error("Failed to get state for " + id + ": " + err.Error() + "\"}")
	} else if valAsBytes == nil {
		return shim.Error("object does not exist: " + id + "\"")
	}
	//unmarshal data to orderHash
	err = json.Unmarshal(valAsBytes, &orderHash)
	if err != nil {
		return shim.Error("cannot unmarshal data")
	}

	hashString := args[1]
	location := args[2]
	status := ""

	//verify hash string
	if orderHash.AssetHash == hashString {
		status = "verify successul"
	} else {
		status = "verify failed"
	}
	ObjectType := "VerifyShipper"
	verify := &VerifyShipper{ObjectType, id, hashString, status, location}
	VerifyToByte, errVerify := json.Marshal(verify)
	if errVerify != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData("verifyShipperCollection", id, VerifyToByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	//create key
	indexName := "OrderID~Hash"
	orderHashIndexKey, errKey := stub.CreateCompositeKey(indexName, []string{ObjectType, id, hashString, status, location})
	if errKey != nil {
		return shim.Error(errKey.Error())
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("verifyShipperCollection", orderHashIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction verifyShipper")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end verifyShipper function ===============")
	return shim.Success(nil)
}

func (t *COD_chaincode) encrypAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start encrypAsset function ===============")
	start := time.Now()
	time.Sleep(time.Second)
	if len(args) != 5 {
		return shim.Error("expecting 5 argument")
	}

	sellerID := args[0]
	asset := args[1]
	detail := args[2]
	quantity, err1 := strconv.Atoi(args[3])
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	price, err2 := strconv.Atoi(args[4])
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	hash := sha256.New()
	hash.Write([]byte(sellerID + asset + detail + string(quantity) + string(price)))
	md := hash.Sum(nil)
	asset_hash := hex.EncodeToString(md)

	fmt.Println("order's hash: ", asset_hash)
	// fmt.Println("image's hash: ", imageHash)

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("time start: ", start)
	fmt.Println("time end: ", end)
	fmt.Println("execute time: ", elapsed)
	printMemUsage()
	fmt.Println("\n=============== end encrypAsset function ===============")

	return shim.Success(nil)
}

func (t *COD_chaincode) dealLimitTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start dealLimitTime function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	if len(args) != 5 {
		return shim.Error("expecting 5 argument")
	}
	orderID := args[0]
	sellerID := args[1]
	deliveryID := args[2]
	orderTime := args[3]
	orderDay := args[4]

	ObjectType := "LimitTime"
	limitTime := &LimitTime{ObjectType, orderID, sellerID, deliveryID, orderTime, orderDay}
	limitTimeToByte, errLimitTime := json.Marshal(limitTime)
	if errLimitTime != nil {
		return shim.Error(errLimitTime.Error())
	}

	errLimitTime = stub.PutPrivateData("limitTimeCollection", orderID, limitTimeToByte)
	if errLimitTime != nil {
		return shim.Error(errLimitTime.Error())
	}

	//create key
	indexName := "orderID~sellerID"
	orderIDIndexKey, errKey := stub.CreateCompositeKey(indexName, []string{ObjectType, orderID, sellerID, deliveryID, orderTime, orderDay})
	if errKey != nil {
		return shim.Error(errKey.Error())
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("limitTimeCollection", orderIDIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Printf("take %s", elapsed.String())
	printMemUsage()
	fmt.Println("\n=============== end dealLimitTime function ===============")

	return shim.Success(nil)
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v MiB", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
