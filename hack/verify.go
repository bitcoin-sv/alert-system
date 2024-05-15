package main

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/bitcoin-sv/alert-system/app/models/model"
)

func main() {
	data, _ := hex.DecodeString("0100000001000000ffcd44660000000001000000115468697320697320616e20616c6572742e1c798b9a1a863f05c8d013bf17d5d15565574acee9dd6276cfa58ac30df442facb09f3cff7b71384f5d7eff1c96dba8bbfab212b40191da99702a6848118521fdd1cdaee843b040bd5977d26021ab84662af1340eead4c7aaa0e4dd97f0875352cbbcd72db4b747e123efa5be3ceefcd988adf1e44c312741119af3610729f2aaa821c76671b76979befb91a54c358f9b153e85a7ff2ec8490d22e7e4dc8c8e576d9c9260026df5b2188d4172407d64c1b27e4fa0c55270a993cbd49098cbdfdc04ba8")
	keys := []string{
		"0233447ada3b75d2c3b98f10d4246458bae53671058c8e71688e4054a24d796d4f",
		"02df63b0401fd7e8eb74e89691c1900dd6f76ccb33a003ad4ac7e9e8ef85eae45f",
		"03dee089fb69671e723bdad999126817dc0d5bc40ac66f4113bb4f510296bf2eb8",
		"03456017726799E8C3C7AE9482F0CDE6B583A5E131299A3B03568C45E72279A331",
		"02408359DF844FDC9E0F19CD96CDF3AD3BFF1F671CE71A65EFEAB042DD23E3C729",
	}
	_appConfig, err := config.LoadDependencies(context.Background(), models.BaseModels, false)
	if err != nil {
		log.Fatalf("error loading configuration: %s", err.Error())
	}
	defer func() {
		_appConfig.CloseAll(context.Background())
	}()
	_appConfig.GenesisKeys = keys
	err = models.CreateGenesisAlert(
		context.Background(), model.WithAllDependencies(_appConfig),
	)
	if err != nil {
		panic(err)
	}
	a, err := models.NewAlertFromBytes(data, model.WithAllDependencies(_appConfig))
	if err != nil {
		panic(err)
	}
	a.Serialize()
	yes, err := a.AreSignaturesValid(context.TODO())
	if err != nil {
		panic(err)
	}
	log.Printf("Verified: %v", yes)
}
