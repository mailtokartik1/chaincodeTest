package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var local string // Entities
	var localval int // Asset holdings
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	local = args[0]
	localval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("localval = %d\n", localval)

	// Write the state to the ledger
	err = stub.PutState(local, []byte(strconv.Itoa(localval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction Invoke
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}
	if function == "PutState" {
		// Deletes an entity from its state
		return t.PutState(stub, args)
	}

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub *shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	local := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(local)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Puts an entry into the state specified number of times
func (t *SimpleChaincode) PutState(stub *shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var err error

	local := args[0]
	var localval int
	localval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	var count int
	count, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for put count")
	}

	// Put the key in the state in ledger
	Avalbytes, err := stub.GetState(local)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	localval, _ = strconv.Atoi(string(Avalbytes))

	// Write the state back to the ledger

	for i := 0; i < count; i++ {
		err = stub.PutState(local, []byte(strconv.Itoa(localval)))
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var local string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	local = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(local)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + local + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + local + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + local + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
