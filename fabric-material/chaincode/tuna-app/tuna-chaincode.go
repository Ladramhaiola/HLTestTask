package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SupplyChainDemo -
type SupplyChainDemo struct {
}

// Tuna structure, with 4 properties.
//Structure tags are used by encoding/json library
//
type Tuna struct {
	ID          string `json:"id"`
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	Location    string `json:"location"`
	Timestamp   string `json:"timestamp"`
	Agreement   string `json:"agreement"`
	Status      string `json:"status"`
}

/*
 * The Init method *
 called when the Smart Contract "tuna-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function
 -- see initLedger()
*/
func (s *SupplyChainDemo) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "tuna-chaincode"
 The app also specifies the specific smart contract function to call with args
*/
func (s *SupplyChainDemo) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger

	switch function {
	case "queryTuna":
		return s.queryTuna(APIstub, args)
	case "queryAllTuna":
		return s.queryAllTuna(APIstub)
	case "changeStatus":
		return s.changeStatus(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "recordTuna":
		return s.recordTuna(APIstub, args)
	}

	return shim.Error("Invalid function name.")
}

/*
 * The queryTuna method *
Used to view the records of one particular tuna
It takes one argument -- the key for the tuna in question
*/
func (s *SupplyChainDemo) queryTuna(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tunaAsBytes, _ := APIstub.GetState(args[0])
	if tunaAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}
	return shim.Success(tunaAsBytes)
}

/*
 * The initLedger method *
Will add test data (10 tuna catches) to our network
*/
func (s *SupplyChainDemo) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	tuna := []Tuna{
		Tuna{ID: "923F", Location: "67.0006, -70.5476", Timestamp: "1504054225", Temperature: "10", Humidity: "24", Agreement: "In"},
		Tuna{ID: "M83T", Location: "91.2395, -49.4594", Timestamp: "1504057825", Temperature: "10", Humidity: "24", Agreement: "In"},
		Tuna{ID: "T012", Location: "58.0148, 59.01391", Timestamp: "1493517025", Temperature: "10", Humidity: "24", Agreement: "In"},
		Tuna{ID: "P490", Location: "-45.0945, 0.7949", Timestamp: "1496105425", Temperature: "10", Humidity: "24", Agreement: "Out"},
		Tuna{ID: "S439", Location: "-107.6043, 19.5003", Timestamp: "1493512301", Temperature: "10", Humidity: "24", Agreement: "Out"},
		Tuna{ID: "J205", Location: "-155.2304, -15.8723", Timestamp: "1494117101", Temperature: "10", Humidity: "24", Agreement: "In"},
		Tuna{ID: "S22L", Location: "103.8842, 22.1277", Timestamp: "1496104301", Temperature: "10", Humidity: "24", Agreement: "In"},
		Tuna{ID: "EI89", Location: "-132.3207, -34.0983", Timestamp: "1485066691", Temperature: "10", Humidity: "24", Agreement: "In"},
		Tuna{ID: "129R", Location: "153.0054, 12.6429", Timestamp: "1485153091", Temperature: "10", Humidity: "24", Agreement: "Out"},
		Tuna{ID: "49W4", Location: "51.9435, 8.2735", Timestamp: "1487745091", Temperature: "10", Humidity: "24", Agreement: "In"},
	}

	i := 0
	for i < len(tuna) {
		fmt.Println("i is ", i)
		tunaAsBytes, _ := json.Marshal(tuna[i])
		APIstub.PutState(strconv.Itoa(i+1), tunaAsBytes)
		fmt.Println("Added", tuna[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordTuna method *
This method takes in 8 arguments (attributes to be saved in the ledger).
*/
func (s *SupplyChainDemo) recordTuna(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	var tuna = Tuna{
		ID:          args[1],
		Temperature: args[2],
		Humidity:    args[3],
		Location:    args[4],
		Timestamp:   args[5],
		Agreement:   args[6],
		Status:      args[7],
	}

	tunaAsBytes, _ := json.Marshal(tuna)
	err := APIstub.PutState(args[0], tunaAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record tuna catch: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllTuna method *
allows for assessing all the records added to the ledger(all tuna catches)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SupplyChainDemo) queryAllTuna(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTuna:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changeTunaHolder method *
The data in the world state can be updated with who has possession.
This function takes in 2 arguments, tuna id and new holder name.
*/
func (s *SupplyChainDemo) changeStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	tunaAsBytes, _ := APIstub.GetState(args[0])
	if tunaAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}

	tuna := Tuna{}
	json.Unmarshal(tunaAsBytes, &tuna)
	// TODO: check if arg satisfy rules
	tuna.Status = args[1]

	tunaAsBytes, _ = json.Marshal(tuna)
	err := APIstub.PutState(args[0], tunaAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change tuna status: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * main function *
calls the Start function
The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SupplyChainDemo))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
