package main

import (
	"encoding/json"
	"fmt"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	 pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Read - read a generic variable from ledger
//
// Shows Off GetState() - reading a key/value from the ledger
//
// Inputs - Array of strings
//  0
//  key
//  "abc"
// 
// Returns - string
// ============================================================================================================================
func read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)           //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes)                  //send it onward
}

// ============================================================================================================================
// Get Asset - get an asset from ledger
// ============================================================================================================================
func get_asset(stub shim.ChaincodeStubInterface, id string) (PharmaAsset, error) {
	var asset PharmaAsset
	assetAsBytes, err := stub.GetState(id)                  //getState retreives a key/value from the ledger
	if err != nil {                                          //this seems to always succeed, even if key didn't exist
		return asset, errors.New("Failed to find asset - " + id)
	}
	json.Unmarshal(assetAsBytes, &asset)                   //un stringify it aka JSON.parse()

	if asset.Id != id {                                     //test if marble is actually here or just nil
		return asset, errors.New("Failed does not exist - " + id)
	}

	return asset, nil 
}

// ============================================================================================================================
// Get history of asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   PharmaAsset   `json:"value"`
	}
	var history []AuditHistory;
	var asset PharmaAsset

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetId := args[0]
	fmt.Printf("- start getHistoryForAsset: %s\n", assetId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(assetId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	
	for resultsIterator.HasNext() {
		keyModification, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		txID := keyModification.TxId
		historicValue := keyModification.Value

		var tx AuditHistory
		tx.TxId = txID                             //copy transaction id over
		json.Unmarshal(historicValue, &asset)     //un stringify it aka JSON.parse()
		if historicValue == nil {                  //asset has been deleted
			var emptyasset PharmaAsset
			tx.Value = emptyasset                 //copy nil asset
		} else {
			json.Unmarshal(historicValue, &asset) //un stringify it aka JSON.parse()
			tx.Value = asset                      //copy asset over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForMarble returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}