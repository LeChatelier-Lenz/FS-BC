/*
Copyright 2022 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

const (
	channelName   = "mychannel"
	chaincodeName = "events"
)

//var now = time.Now()

//var assetID = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

type Currency struct {
	CurrencyID string  `json:"CurrencyID"` //格式为"Currency"+时间戳
	Amount     float32 `json:"Amount"`     //限制货币的最小单位为0.01
	Owner      string  `json:"Owner"`      //user_id
	CreatedAt  string  `json:"CreatedAt"`
	CreatedVia string  `json:"CreatedVia"` //"Loan","Insurance","Transfer","Deposit","System"
	UpdatedAt  string  `json:"UpdatedAt"`
	UpdatedVia string  `json:"UpdatedVia"` //"Loan","Insurance","Transfer"
}

type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"Color"`
	Size           string `json:"Size"`
	Owner          string `json:"Owner"`
	AppraisedValue string `json:"AppraisedValue"`
}

type QueryAssetRequest struct {
	user_id  string `json:"user_id"`
	password string `json:"password"`
}

type ContractQueryRequest struct {
	user_id  string `json:"user_id"`
	password string `json:"password"`
}
type ContractQueryByIdRequest struct {
	user_id       string `json:"user_id"`
	business_id   string `json:"business_id"`
	business_type string `json:"business_type"`
}
type CreateContractRequest struct {
	user_id        string  `json:"user_id"`
	password       string  `json:"password"`
	current_time   string  `json:"current_time"`
	bussiness_type string  `json:"bussiness_type"`
	amount         float32 `json:"amount"`
	rate           float32 `json:"rate"`
	issuer         string  `json:"issuer"`
	period         int     `json:"period"`
	business_id    string  `json:"business_id"`
}
type Conditions struct {
	credit           float32 `json:"credit"`
	income           float32 `json:"income"`
	is_sudden        bool    `json:"is_sudden"`
	contingency_info string  `json:"contingency_info"`
}

type LoanStartRequest struct {
	user_id      string     `json:"user_id"`
	bussiness_id string     `json:"bussiness_id"`
	conditions   Conditions `json:"conditions"`
	current_time string     `json:"current_time"`
}
type LoanCheckRequest struct {
	user_id      string     `json:"user_id"`
	bussiness_id string     `json:"bussiness_id"`
	conditions   Conditions `json:"conditions"`
	current_time string     `json:"current_time"`
}
type InsuranceStartRequest struct {
	user_id      string     `json:"user_id"`
	bussiness_id string     `json:"bussiness_id"`
	conditions   Conditions `json:"conditions"`
	current_time string     `json:"current_time"`
}
type InsuranceCheckRequest struct {
	user_id      string     `json:"user_id"`
	bussiness_id string     `json:"bussiness_id"`
	conditions   Conditions `json:"conditions"`
	current_time string     `json:"current_time"`
}
type PayTranserRequest struct {
	user_id        string  `json:"user_id"`
	password       string  `json:"password"`
	amount         float32 `json:"amount"`
	target_user_id string  `json:"target_user_id"`
	current_time   string  `json:"current_time"`
}
type DepositTranserRequest struct {
	user_id      string  `json:"user_id"`
	password     string  `json:"password"`
	amount       float32 `json:"amount"`
	current_time string  `json:"current_time"`
}

func main() {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Context used for event listening
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for events emitted by subsequent transactions
	startChaincodeEventListening(ctx, network)

	//调用链码进行资产初始化

	router := gin.Default()

	router.GET("/ecosys/asset", func(c *gin.Context) {

		var queryAssetRequest QueryAssetRequest
		if err := c.ShouldBindJSON(&queryAssetRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

		evaluateResult, err := contract.EvaluateTransaction("ReadTotalCurrencyByOwner", queryAssetRequest.user_id, queryAssetRequest.password)
		if err != nil {
			panic(fmt.Errorf("failed to evaluate transaction: %w", err))
		}
		result := formatJSON(evaluateResult)

		fmt.Printf("*** Result:%s\n", result)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "GetAllAssets Failed",
				"result":  "",
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "GetAllAssets Success",
			"result":  result,
		})
	})

	router.GET("/ecosys/contract", func(c *gin.Context) {
		var contractQueryRequest ContractQueryRequest
		if err := c.ShouldBindJSON(&contractQueryRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("\n--> Evaluate Transaction: GetAllContracts, function returns all the current Contracts on the ledger")

		insuranceResult, err1 := contract.EvaluateTransaction("ReadInsuranceListByOwner", contractQueryRequest.user_id)
		if err1 != nil {
			panic(fmt.Errorf("failed to evaluate transaction: %w", err1))
		}

		loanResult, err2 := contract.EvaluateTransaction("ReadLoanListByOwner")
		if err2 != nil {
			panic(fmt.Errorf("failed to evaluate transaction: %w", err2))
		}
		// 把insuranceResult 和 loanResult 合并
		evaluateResult := bytes.Join([][]byte{insuranceResult, loanResult}, []byte("\n"))

		result := formatJSON(evaluateResult)

		fmt.Printf("*** Result:%s\n", result)

		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "GetAllContracts Failed",
				"result":  "",
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err1)
			return
		}
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "GetAllContracts Failed",
				"result":  "",
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err2)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "GetAllContracts Success",
			"result":  result,
		})
	})
	router.GET("/ecosys/query_contract", func(c *gin.Context) {
		var contractQueryByIdRequest ContractQueryByIdRequest
		if err := c.ShouldBindJSON(&contractQueryByIdRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("\n--> Evaluate Transaction: GetAllContracts, function returns all the current Contracts on the ledger")

		if contractQueryByIdRequest.business_type == "loan" {
			evaluateResult, err := contract.EvaluateTransaction("ReadLoan", contractQueryByIdRequest.user_id, contractQueryByIdRequest.business_id)
			result := formatJSON(evaluateResult)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    "500",
					"message": "GetAllContracts Failed",
					"result":  "",
				})
				log.Fatalf("Failed to evaluate transaction: %s\n", err)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    "200",
				"message": "GetAllContracts Success",
				"result":  result,
			})
		} else {
			evaluateResult, err := contract.EvaluateTransaction("ReadInsurance", contractQueryByIdRequest.user_id, contractQueryByIdRequest.business_id)
			result := formatJSON(evaluateResult)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    "500",
					"message": "GetAllContracts Failed",
					"result":  "",
				})
				log.Fatalf("Failed to evaluate transaction: %s\n", err)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    "200",
				"message": "GetAllContracts Success",
				"result":  result,
			})
		}
	})

	router.POST("/ecosys/new_contract", func(c *gin.Context) {
		var createContractRequest CreateContractRequest
		if err := c.ShouldBindJSON(&createContractRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("\n--> Evaluate Transaction: Create a new contract, function returns all the current Contracts on the ledger")
		evaluateResult, err := contract.SubmitTransaction("CreateContract", createContractRequest.user_id, createContractRequest.business_id, fmt.Sprintf("%f", createContractRequest.amount), createContractRequest.issuer, fmt.Sprintf("%f", createContractRequest.rate), createContractRequest.bussiness_type, fmt.Sprintf("%d", createContractRequest.period))
		result := formatJSON(evaluateResult)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Create Failed",
				"result":  string(result),
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Create Success",
			"result":  string(result),
		})
	})

	//firstBlockNumber := createAsset(contract)

	router.POST("/ecosys/loan/start", func(c *gin.Context) {
		var loanStartRequest LoanStartRequest
		err := c.BindJSON(&loanStartRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Bad Request",
				"result":  "",
			})
			return
		}
		result, err := contract.SubmitTransaction("StartLoan", loanStartRequest.user_id, loanStartRequest.bussiness_id, fmt.Sprintf("%f", loanStartRequest.conditions.credit), fmt.Sprintf("%f", loanStartRequest.conditions.income))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Update Failed",
				"result":  string(result),
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Update Success",
			"result":  string(result),
		})
	})

	router.POST("/ecosys/loan/check_contract", func(c *gin.Context) {
		var loanCheckRequest LoanCheckRequest
		err := c.BindJSON(&loanCheckRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Bad Request",
				"result":  "",
			})
			return
		}
		result, err := contract.SubmitTransaction("LoanContractCheck", loanCheckRequest.user_id, loanCheckRequest.bussiness_id, fmt.Sprintf("%f", loanCheckRequest.conditions.credit), fmt.Sprintf("%f", loanCheckRequest.conditions.income))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Loan Check Failed",
				"result":  result,
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Loan Check Success",
			"result":  result,
		})
	})

	router.POST("/ecosys/insurance/start", func(c *gin.Context) {
		var insuranceStartRequest InsuranceStartRequest
		err := c.BindJSON(&insuranceStartRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Bad Request",
				"result":  "",
			})
			return
		}
		result, err := contract.SubmitTransaction("StartInsurance", insuranceStartRequest.user_id, insuranceStartRequest.bussiness_id, fmt.Sprintf("%f", insuranceStartRequest.conditions.credit), fmt.Sprintf("%f", insuranceStartRequest.conditions.income))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Insurance Start Failed",
				"result":  string(result),
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Insurance Start Success",
			"result":  string(result),
		})
	})
	router.POST("/ecosys/insurance/check_contract", func(c *gin.Context) {
		var insuranceCheckRequest InsuranceCheckRequest
		err := c.BindJSON(&insuranceCheckRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Bad Request",
				"result":  "",
			})
			return
		}
		result, err := contract.SubmitTransaction("InsuranceContractCheck", insuranceCheckRequest.user_id, insuranceCheckRequest.bussiness_id, fmt.Sprintf("%f", insuranceCheckRequest.conditions.credit), fmt.Sprintf("%f", insuranceCheckRequest.conditions.income), fmt.Sprintf("%t", insuranceCheckRequest.conditions.is_sudden), insuranceCheckRequest.conditions.contingency_info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Insurance Check Failed",
				"result":  result,
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Insurance Check Success",
			"result":  result,
		})
	})
	router.POST("/ecosys/pay/transfer", func(c *gin.Context) {
		var payTransferRequest PayTranserRequest
		err := c.BindJSON(&payTransferRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Bad Request",
				"result":  "",
			})
			return
		}
		result, err := contract.SubmitTransaction("TransferCurrency", payTransferRequest.user_id, payTransferRequest.target_user_id, fmt.Sprintf("%f", payTransferRequest.amount), "Transfer")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Pay Transfer Failed",
				"result":  result,
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Pay Transfer Success",
			"result":  result,
		})
	})
	router.POST("/ecosys/pay/deposit", func(c *gin.Context) {
		var depositTransferRequest DepositTranserRequest
		err := c.BindJSON(&depositTransferRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Bad Request",
				"result":  "",
			})
			return
		}
		var currency = Currency{
			CurrencyID: "Currency" + depositTransferRequest.current_time,
			Amount:     depositTransferRequest.amount,
			Owner:      depositTransferRequest.user_id,
			CreatedAt:  depositTransferRequest.current_time,
			CreatedVia: "Deposit",
			UpdatedAt:  depositTransferRequest.current_time,
			UpdatedVia: "Deposit",
		}
		currencyJson, err := json.Marshal(currency)
		currencyJsonStr := string(currencyJson)
		result, err := contract.SubmitTransaction("CreateCurrency", currencyJsonStr) // TODO
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Deposite Failed",
				"result":  result,
			})
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Deposite Success",
			"result":  result,
		})
	})

	err = router.Run(":8000")
	if err != nil {
		return
	}

	// Replay events from the block containing the first transaction
	//replayChaincodeEvents(ctx, network, firstBlockNumber)
}

func startChaincodeEventListening(ctx context.Context, network *client.Network) {
	fmt.Println("\n*** Start chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		for event := range events {
			asset := formatJSON(event.Payload)
			fmt.Printf("\n<-- Chaincode event received: %s - %s\n", event.EventName, asset)
		}
	}()
}

func formatJSON(data []byte) string {
	var result bytes.Buffer
	if err := json.Indent(&result, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return result.String()
}

func createAsset(contract *client.Contract, asset Asset) ([]byte, error) {
	var err error
	fmt.Printf("\n--> Submit transaction: CreateAsset, %s owned by %s with appraised value %s\n", asset.ID, asset.Owner, asset.AppraisedValue)

	result, commit, err := contract.SubmitAsync("CreateAsset", client.WithArguments(
		asset.ID,
		asset.Color,
		asset.Size,
		asset.Owner,
		asset.AppraisedValue,
	))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	status, err := commit.Status()
	if err != nil {
		panic(fmt.Errorf("failed to get transaction commit status: %w", err))
	}

	if !status.Successful {
		panic(fmt.Errorf("failed to commit transaction with status code %v", status.Code))
	}
	fmt.Println("\n*** CreateAsset committed successfully")

	if err != nil {
		return result, err
	}
	return result, nil
}

func updateAsset(contract *client.Contract, asset Asset) ([]byte, error) {
	var err error
	fmt.Printf("\n--> Submit transaction: UpdateAsset, %s update appraised value to 200\n", asset.ID)

	result, err := contract.SubmitTransaction("UpdateAsset", asset.ID, asset.Color, asset.Size, asset.Owner, asset.AppraisedValue)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Println("\n*** UpdateAsset committed successfully")

	if err != nil {
		return result, err
	}
	return result, nil
}

func transferAsset(contract *client.Contract, assetID string, newOwner string) error {
	fmt.Printf("\n--> Submit transaction: TransferAsset, %s to Mary\n", assetID)

	_, err := contract.SubmitTransaction("TransferAsset", assetID, newOwner)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Println("\n*** TransferAsset committed successfully")
	if err != nil {
		return err
	}
	return nil
}

func deleteAsset(contract *client.Contract, assetID string) error {
	fmt.Printf("\n--> Submit transaction: DeleteAsset, %s\n", assetID)

	_, err := contract.SubmitTransaction("DeleteAsset", assetID)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Println("\n*** DeleteAsset committed successfully")
	if err != nil {
		return err
	}
	return nil
}
