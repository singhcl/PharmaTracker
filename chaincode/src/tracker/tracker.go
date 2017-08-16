package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type PharmaTrackerChaincode struct {
}

// ============================================================================================================================
// Asset Definitions - The ledger will store assets
// ============================================================================================================================

// ----- Marbles ----- //
type PharmaAsset struct {
	Id         		string        		`json:"id"`      
	Type       		string        		`json:"type"`
	Category   		string        		`json:"category"`
	AssetClass      string        		`json:"assetClass"`
	Data       		AssetData     		`json:"assetData"`    
	TraceData  		AssetTraceData 		`json:"assetTraceData"`
}

type AssetData struct {
	Info         AssetInfo 		 `json:"info"`
	Children     []AssetChildren `json:"children"`    
}

type AssetTraceData struct {
	Owner         		string `json:"owner"`
	Status   		 	string `json:"status"`
	EventDateTime      	string `json:"eventDateTime"`
	Location         	string `json:"location"`
	GeoLocation   		string `json:"geoLocation"`
}

type AssetInfo struct {
	Name         	string `json:"name"`
	Type   		 	string `json:"type"`
	PkgSize      	int    `json:"pkgSize"`
	MfgDate         string `json:"mfgDate"`
	LotNo   		string `json:"lotNo"`
	ExpiryDate      string `json:"expiryDate"`
}

type AssetChildren struct {
	Id         string 	`json:"id"`
	Type       string 	`json:"type"`    
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(PharmaTrackerChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}


// ============================================================================================================================
// Init - initialize the chaincode - No initialization required
// ============================================================================================================================
func (t *PharmaTrackerChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("PharmaTrackerChaincode Is Starting Up")
	return shim.Success(nil)
}


// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *PharmaTrackerChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "init" {                    //initialize the chaincode state, used as reset
		return t.Init(stub)
	} else if function == "fetch" {             //generic read ledger
		return get_asset(stub, args)
	} else if function == "write" {            //generic writes to ledger
		return write_asset(stub, args)
	} else if function == "delete" {           //deletes an asset from state
		return delete_asset(stub, args)
	} else if function == "getHistory"{        //read history of an asset (audit)
		return getHistory(stub, args)
	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}


// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *PharmaTrackerChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}
