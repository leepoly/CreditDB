package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions
type SmartContract struct {
	contractapi.Contract
}

type Loan struct {
	Value      string `json:"value"`
	SenderName string `json:"sendername"`
	RecverName string `json:"recvername"`
	Timestamp  string `json:"timestamp"`
}

// QueryResult structure used for handling result of query
type TxQueryResult struct {
	Key    string `json:"Key"`
	Record *Loan
}

// InitLedger adds a base set of loans to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	loans := []Loan{
		Loan{Value: "5.00", SenderName: "Bank of Beijing", RecverName: "YiweiLi", Timestamp: "2020-09-02 15:04:05"},
		Loan{Value: "3.00", SenderName: "YiweiLi", RecverName: "Bank of Beijing", Timestamp: "2020-10-03 15:04:05"},
		Loan{Value: "3.00", SenderName: "Bank of China", RecverName: "Test", Timestamp: "2020-11-01 15:59:59"},
		Loan{Value: "3.00", SenderName: "Test", RecverName: "Bank of China", Timestamp: "2020-11-02 15:59:59"},
	}

	for i, loan := range loans {
		loanAsBytes, _ := json.Marshal(loan)
		err := ctx.GetStub().PutState("TX"+strconv.Itoa(i), loanAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateLoan adds a new Tx to the world state with given details
func (s *SmartContract) CreateLoan(ctx contractapi.TransactionContextInterface, ID string, Value string, SenderName string, RecverName string, Timestamp string) error {
	loan := Loan{
		Value:      Value,
		SenderName: SenderName,
		RecverName: RecverName,
		Timestamp:  Timestamp,
	}

	loanAsBytes, _ := json.Marshal(loan)

	return ctx.GetStub().PutState("TX"+ID, loanAsBytes)
}

// QueryLoan returns the loan stored in the world state with given id
func (s *SmartContract) QueryLoan(ctx contractapi.TransactionContextInterface, loanNumber string) (*Loan, error) {
	loanNumber = "TX" + loanNumber
	loanAsBytes, err := ctx.GetStub().GetState(loanNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if loanAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", loanNumber)
	}

	loan := new(Loan)
	_ = json.Unmarshal(loanAsBytes, loan)

	return loan, nil
}

// ListLoans returns all loans found in world state
func (s *SmartContract) ListLoans(ctx contractapi.TransactionContextInterface) ([]TxQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []TxQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		loan := new(Loan)
		_ = json.Unmarshal(queryResponse.Value, loan)

		queryResult := TxQueryResult{Key: queryResponse.Key, Record: loan}
		results = append(results, queryResult)
	}

	return results, nil
}

// Queryuser returns all loans related to a username (sender or receiver)
func (s *SmartContract) QueryUser(ctx contractapi.TransactionContextInterface, QueryName string) ([]TxQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []TxQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		loan := new(Loan)
		_ = json.Unmarshal(queryResponse.Value, loan)

		if loan.RecverName == QueryName || loan.SenderName == QueryName {
			queryResult := TxQueryResult{Key: queryResponse.Key, Record: loan}
			results = append(results, queryResult)
		}
	}

	return results, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create creditdb chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting creditdb chaincode: %s", err.Error())
	}
}
