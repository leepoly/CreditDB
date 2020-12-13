package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath" // General packages
	"strconv"
	"time"

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

	contract := network.GetContract("creditdb")
	return contract
}

func listTx(contract *gateway.Contract) ([]map[string]interface{}, error) {
	result, contractErr := contract.EvaluateTransaction("listLoans")
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

func queryUser(contract *gateway.Contract, username string) ([]map[string]interface{}, error) {
	result, contractErr := contract.EvaluateTransaction("QueryUser", username)
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

func queryTx(contract *gateway.Contract, ID string) (map[string]interface{}, error) {
	result, contractErr := contract.EvaluateTransaction("queryLoan", ID)
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

func createTx(contract *gateway.Contract, ID string, Value float64, SenderName string, RecverName string, currentTime string) string {
	ValueStr := fmt.Sprintf("%.2f", Value)
	result, contractErr := contract.SubmitTransaction("CreateLoan", ID, ValueStr, SenderName, RecverName, currentTime)
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
			"Introduion":  "A digital credit system based on hyperledger",
			"Description": "THU Cyber Intellective Economic and Block Chain, Course project",
			"RepoLink":    "https://github.com/leepoly/CreditDB",
			"Version":     "0.1",
		})
	})

	return router
}

func main() {
	contract := initWallet()
	r := initServer()

	r.GET("/listTx", func(c *gin.Context) {
		jsonResult, err := listTx(contract)
		if err == nil {
			c.IndentedJSON(http.StatusOK, jsonResult)
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{
				"Err": err.Error(),
			})
		}
	})

	r.POST("/queryUser", func(c *gin.Context) {
		username := c.PostForm("name")
		jsonResult, err := queryUser(contract, username)
		if err == nil {
			c.IndentedJSON(http.StatusOK, jsonResult)
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{
				"Err": err.Error(),
			})
		}
	})

	r.POST("/queryTx", func(c *gin.Context) {
		id := c.PostForm("id")
		jsonResult, err := queryTx(contract, id)
		if err == nil {
			c.IndentedJSON(http.StatusOK, jsonResult)
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{
				"Err": err.Error(),
			})
		}
	})

	r.POST("/createTx", func(c *gin.Context) {
		id := c.PostForm("id")
		valueStr := c.PostForm("value")
		valueFloat64, _ := strconv.ParseFloat(valueStr, 32)
		SenderName := c.DefaultPostForm("SenderName", "defaultsender")
		RecverName := c.DefaultPostForm("Recvername", "defaultrecver")
		currentTime := time.Now().Format("01-02-2006 15:04:05")
		result := createTx(contract, id, valueFloat64, SenderName, RecverName, currentTime)
		c.IndentedJSON(http.StatusOK, gin.H{
			"Err": result,
		})
	})

	r.Run(":8085") // launch the gin engine
}
