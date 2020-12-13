package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath" // General packages

	"github.com/gin-gonic/gin" // Gin backend

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway" // Hyperledger Go SDK
)

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}

func initWallet() *gateway.Contract {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("fabcar")
	return contract
}

func listAccount(contract *gateway.Contract) ([]map[string]interface{}, error) {
	result, contractErr := contract.EvaluateTransaction("queryAllCars")
	if contractErr != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", contractErr)
		return nil, contractErr
	}

	var jsonResult []map[string]interface{}
	jsonErr := json.Unmarshal([]byte(result), &jsonResult)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}
	return jsonResult, jsonErr
}

func queryAccount(contract *gateway.Contract, name string) (map[string]interface{}, error) {
	result, contractErr := contract.EvaluateTransaction("queryCar", name)
	if contractErr != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", contractErr)
		return nil, contractErr
	}

	var jsonResult map[string]interface{}
	jsonErr := json.Unmarshal([]byte(result), &jsonResult)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}
	return jsonResult, jsonErr
}

func createAccount(contract *gateway.Contract, name string, field1 string, field2 string, field3 string, field4 string) string {
	result, contractErr := contract.SubmitTransaction("createCar", name, field1, field2, field3, field4)
	if contractErr != nil {
		fmt.Printf("Failed to submit transaction: %s\n", contractErr)
		return contractErr.Error()
	}
	return string(result[:])
}

func modifyAccountField4(contract *gateway.Contract, name string, field4NewValue string) string {
	result, contractErr := contract.SubmitTransaction("changeCarOwner", name, field4NewValue)
	if contractErr != nil {
		fmt.Printf("Failed to submit transaction: %s\n", contractErr)
		return contractErr.Error()
	}
	return string(result[:])
}

func initServer() *gin.Engine {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"Introduion": "A digital credit system based on hyperledger",
			"Authors":    "THU Cyber Intellective Economic and Block Chain, Course project",
			"Version":    "0.1",
		})
	})

	return router
}

func main() {
	contract := initWallet()
	r := initServer()

	r.GET("/listAccount", func(c *gin.Context) {
		jsonResult, err := listAccount(contract)
		if err == nil {
			c.IndentedJSON(http.StatusOK, jsonResult)
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{
				"Err": err.Error(),
			})
		}
	})

	r.POST("/queryAccount", func(c *gin.Context) {
		name := c.PostForm("name")
		jsonResult, err := queryAccount(contract, name)
		if err == nil {
			c.IndentedJSON(http.StatusOK, jsonResult)
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{
				"Err": err.Error(),
			})
		}
	})

	r.POST("/createAccount", func(c *gin.Context) {
		name := c.PostForm("name")
		field1 := c.DefaultPostForm("field1", "VW")
		field2 := c.DefaultPostForm("field2", "Polo")
		field3 := c.DefaultPostForm("field3", "Grey")
		field4 := c.DefaultPostForm("field4", "Mary")
		result := createAccount(contract, name, field1, field2, field3, field4)
		c.IndentedJSON(http.StatusOK, gin.H{
			"Err": result,
		})
	})

	r.POST("/modifyAccountField4", func(c *gin.Context) {
		name := c.PostForm("name")
		field4 := c.DefaultPostForm("field4", "Archie")
		result := modifyAccountField4(contract, name, field4)
		c.IndentedJSON(http.StatusOK, gin.H{
			"Err": result,
		})
	})

	r.Run(":8085") // launch the gin engine
}
