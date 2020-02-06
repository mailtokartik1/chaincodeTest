/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init initializes the chaincode
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println("SampleChainCode Init")

	//
	// Demonstrate the use of Attribute-Based Access Control (ABAC) by checking
	// to see if the caller has the "abac.init" attribute with a value of true;
	// if not, return an error.
	//
	err := cid.AssertAttributeValue(stub, "sampleChainCode.init", "true")
	if err != nil {
		return shim.Error(err.Error())
	}

	_, args := stub.GetFunctionAndParameters()
	var local string // Entities
	var localVal int // Asset holdings

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	local = args[0]
	localVal, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	fmt.Printf("localVal = %d", localVal)

	// Write the state to the ledger
	err = stub.PutState(local, []byte(strconv.Itoa(localVal)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Chaincode invoked")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "delete":
		return t.delete(stub, args)
	case "query":
		return t.query(stub, args)
	case "putMultiple":
		return t.putMultiple(stub, args)
	default:
		return shim.Error(fmt.Sprintf("Invalid Smart Contract function : %s", function))
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\" \"putMultiple\"")
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	local := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(local)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("Deleted")

	return shim.Success(nil)
}

// Puts an entry into the state specified number of times
func (t *SimpleChaincode) putMultiple(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	var err error
	local := args[0]
	var localVal, times int
	localVal, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	times, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	for i := 0; i < times; i++ {
		err = stub.PutState(local, []byte(strconv.Itoa(localVal)))
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	fmt.Println("State put %d times", times)

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var local string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	local = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(local)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + local + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + local + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + local + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
