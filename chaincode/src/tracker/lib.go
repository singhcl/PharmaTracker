package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

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
		return asset, errors.New("Asset does not exist - " + id)
	}

	return asset, nil
}

// ========================================================
// Input Sanitation - input checking, look for required checks
// ========================================================
func sanitize_arguments(strs []string) error{
	for i, val:= range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(val) > 32 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
		}
	}
	return nil
}
