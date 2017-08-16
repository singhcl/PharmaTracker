package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// delete_asset() - remove a asset from state and from asset index
// 
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs - Array of strings
//      0      ,         1
//     id      ,  authed_by_company
// "m999999999", "united assets"
// ============================================================================================================================
func delete_asset(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete_asset")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	id := args[0]
	// get the asset
	asset, err := get_asset(stub, id)
	if err != nil{
		fmt.Println("Failed to find asset by id " + id)
		return shim.Error(err.Error())
	}

	// remove the asset
	err = stub.DelState(id)                                                 //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_asset")
	return shim.Success(nil)
}

// ============================================================================================================================
// Write PharmaAsset - create a new asset, store into chaincode state
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings
//      0      ,    1  ,  2  ,      3          ,       4
//     id      ,  color, size,     owner id    ,  authing company
// "m999999999", "blue", "35", "o9999999999999", "united assets"
// ============================================================================================================================
func write_asset(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting init_asset")

	/*if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}*/

	id := args[0]
	color := strings.ToLower(args[1])
	owner_id := args[3]
	authed_by_company := args[4]
	size, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}

	//check if asset id already exists
	asset, err := get_asset(stub, id)
	if err == nil {
		fmt.Println("This asset already exists - " + id)
		fmt.Println(asset)
		return shim.Error("This asset already exists - " + id)  //all stop a asset by this id exists
	}

	//build the asset json string manually
	str := `{
		"docType":"asset", 
		"id": "` + id + `", 
		"color": "` + color + `", 
		"size": ` + strconv.Itoa(size) + `, 
		"owner": {
			"id": "` + owner_id + `", 
			"username": "` + owner.Username + `", 
			"company": "` + owner.Company + `"
		}
	}`
	err = stub.PutState(id, []byte(str))                         //store asset with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_asset")
	return shim.Success(nil)
}