package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func checkInit(t *testing.T, stub *shimtest.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shimtest.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkUpdate(t *testing.T, stub *shimtest.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("update"), []byte(name), []byte(value)})
	if res.Status != shim.OK {
		fmt.Println("Update", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Update", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Update value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shimtest.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("query"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkPutMultiple(t *testing.T, stub *shimtest.MockStub, name string, value string, times string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("putMultiple"), []byte(name), []byte(value), []byte(times)})
	if res.Status != shim.OK {
		fmt.Println("Putting multiple times", name, "failed", string(res.Message))
		t.FailNow()
	}
}

func TestChaincodeInit(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shimtest.NewMockStub("chaincode", scc)

	// Init

	checkInit(t, stub, [][]byte{[]byte("init"), []byte("A"), []byte("123")})

	// Check state

	checkState(t, stub, "A", "123")

	// Init

	checkInit(t, stub, [][]byte{[]byte("init"), []byte("B"), []byte("456")})

	// Check state

	checkState(t, stub, "B", "456")
}

func TestChaincodeQuery(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shimtest.NewMockStub("abac", scc)

	// Init

	checkInit(t, stub, [][]byte{[]byte("init"), []byte("A"), []byte("345")})

	// Query for expected value

	checkQuery(t, stub, "A", "345")

	// Init

	checkInit(t, stub, [][]byte{[]byte("init"), []byte("B"), []byte("456")})

	// Query for expected value

	checkQuery(t, stub, "B", "456")
}

func TestChaincodeUpdate(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shimtest.NewMockStub("abac", scc)

	// Init

	checkInit(t, stub, [][]byte{[]byte("init"), []byte("A"), []byte("345")})

	// Check updated value

	checkUpdate(t, stub, "A", "345")

	// Init

	checkInit(t, stub, [][]byte{[]byte("init"), []byte("B"), []byte("456")})

	// Check updated value

	checkUpdate(t, stub, "B", "456")
}

func TestChaincodePutMultiple(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shimtest.NewMockStub("abac", scc)

	// Init

	checkInit(t, stub, [][]byte{[]byte("init"), []byte("A"), []byte("345")})

	// Check final output

	checkPutMultiple(t, stub, "A", "345", "10")

	// Query final value

	checkQuery(t, stub, "A", "345")
}
