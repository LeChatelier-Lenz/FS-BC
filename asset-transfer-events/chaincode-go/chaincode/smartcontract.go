package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

/* 金融履约系统链码介绍 */
// 金融履约系统链码是一个简单的金融履约系统，包含了货币、保险、贷款三个基本功能。
// 该链码是一个简单的示例，用于展示金融履约系统的基本功能，包括货币的创建、转移，保险合同的创建、启动、赔偿，贷款合同的创建、启动、强制还款等。
// 注意事项：
// 1.关于货币的创建和转移，使用了UTXO方式，即每笔交易都是一个新的货币，不会出现找零的情况。
// 2.关于保险合同的启动和赔偿，需要根据申请人的信用分和收入情况进行判断，以决定是否启动保险合同，以及是否需要赔偿。
// 3.关于贷款合同的启动和强制还款，需要根据申请人的信用分和收入情况进行判断，以决定是否启动贷款合同，以及是否需要强制还款。
// 4.关于区块链中的“键”，使用了复合键的方式，即将多个键（owner+id）组合在一起，作为一个复合键，用于查询资产。因此在查询资产时，需要注意输入。而且id本身又是一个自带“资产名”+时间戳的特殊形式，需要注意。
// 5.合同调用链码全流程：
//    step1: 客户端创建合同时，实际调用链码的Create函数，创建合同。Create函数会创建一个“Applied”状态的合同。 - CreateContract
//    step2: 调用Start函数，启动合同，支付保险金/贷款金额。 - StartLoan/StartInsurance
//    step2: 合同启动后，根据合同的状态，进行后续操作，如保险合同的赔偿，贷款合同的强制还款等。 - InsuranceContractCheck/LoanContractCheck
//    step3: 合同的状态变化，会触发相应的事件，客户端可以监听事件，进行后续操作。
// 6.资产查询调用链码全流程：
//    ReadTotalCurrencyByOwner 查询某个用户的当前总余额
//    ReadLoanListByOwner 通过owner查询贷款合同列表
//    ReadInsuranceListByOwner 通过owner查询保险合同列表
// 7.支付行为调用链码全流程：
//    TransferCurrency 货币结构体的转移函数，使用UTXO方式。该函数体现了货币的使用方式，即转账。（注意，不再使用合同方式操作了）

/* Currency 全流程
 * 货币结构体，作为交易其他资产的基础，可以被转让，用来作为系统中用户的账户余额
 * CreateCurrency 货币结构体的创建函数，用于创建系统货币/用户存入货币。
 * ReadCurrency 根据id读取货币
 * ReadCurrencyListByOwner 通过owner查询货币列表，是一个辅助函数
 * ReadTotalCurrencyByOwner 查询某个用户（owner）的当前总余额
 * TransferCurrency 货币结构体的转移函数，使用UTXO方式。该函数体现了货币的使用方式，即转账。
 */

// Currency 系统货币结构体，作为交易其他资产的基础，可以被转让，用来作为系统中用户的账户余额
type Currency struct {
	CurrencyID string  `json:"CurrencyID"` //格式为"Currency"+时间戳
	Amount     float32 `json:"Amount"`     //限制货币的最小单位为0.01
	Owner      string  `json:"Owner"`      //user_id
	CreatedAt  string  `json:"CreatedAt"`
	CreatedVia string  `json:"CreatedVia"` //"Loan","Insurance","Transfer","Deposit","System"
	UpdatedAt  string  `json:"UpdatedAt"`
	UpdatedVia string  `json:"UpdatedVia"` //"Loan","Insurance","Transfer"
}

// CreateCurrency 货币结构体的创建函数，用于创建系统货币/用户存入货币。
// currencyBytes 参数是一个json格式的货币结构体(需要先转换为[]byte)
func (s *SmartContract) CreateCurrency(ctx contractapi.TransactionContextInterface, currencyBytes []byte) error {
	var currency Currency
	err := json.Unmarshal(currencyBytes, &currency)
	// 检查货币是否已经存在
	compositeKey, err := ctx.GetStub().CreateCompositeKey("Currency", []string{currency.Owner, currency.CurrencyID})
	existing, err := s.readState(ctx, compositeKey)
	if err == nil && existing != nil {
		return fmt.Errorf("the asset %s already exists", currency.CurrencyID)
	}
	assetJSON, err := json.Marshal(currency)
	if err != nil {
		return err
	}
	ctx.GetStub().SetEvent("CreateCurrency", assetJSON)
	return ctx.GetStub().PutState(compositeKey, assetJSON)
}

// ReadCurrency 根据id读取货币
// 因为是直接调用id，所以是一个链码内部函数，不需要暴露给外部
func (s *SmartContract) ReadCurrency(ctx contractapi.TransactionContextInterface, id string) (*Currency, error) {
	assetJSON, err := s.readState(ctx, id)
	if err != nil {
		return nil, err
	}

	var asset Currency
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// ReadCurrencyListByOwner 通过owner查询货币列表，是一个辅助函数
func (s *SmartContract) ReadCurrencyListByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]Currency, error) {
	// 通过owner查询货币列表
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Currency", []string{owner})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var currencyList []Currency
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var currency Currency
		err = json.Unmarshal(queryResponse.Value, &currency)
		if err != nil {
			return nil, err
		}
		if currency.Owner == owner {
			currencyList = append(currencyList, currency)
		}
	}
	return currencyList, nil
}

// ReadTotalCurrencyByOwner 查询某个用户（owner）的当前总余额
func (s *SmartContract) ReadTotalCurrencyByOwner(ctx contractapi.TransactionContextInterface, owner string) (float32, error) {
	currencyList, err := s.ReadCurrencyListByOwner(ctx, owner)
	if err != nil {
		return 0, err
	}
	var totalAmount float32
	for _, currency := range currencyList {
		totalAmount += currency.Amount
	}
	return totalAmount, nil
}

// TransferCurrency 货币结构体的转移函数，使用UTXO方式。该函数体现了货币的使用方式，即转账。
// transferReason是转账原因，可以是"Loan","Insurance","Transfer"等,用于记录货币的使用情况。也供函数调用时指明转账原因。
func (s *SmartContract) TransferCurrency(ctx contractapi.TransactionContextInterface, oldOwner string, newOwner string, amount float32, transferReason string) error {
	oldCurrencyList, err := s.ReadCurrencyListByOwner(ctx, oldOwner)
	if err != nil {
		log.Fatalln("Failed to read currency list by owner")
		return err
	}
	if len(oldCurrencyList) == 0 {
		return fmt.Errorf("no currency found for owner %s", oldOwner)
	}
	var totalAmount float32
	var DeleteCurrencyList []Currency
	// 遍历货币列表，找到足够的货币转账
	for _, currency := range oldCurrencyList {
		totalAmount += currency.Amount
		DeleteCurrencyList = append(DeleteCurrencyList, currency)
		if totalAmount >= amount {
			break
		}
	}
	// 检查余额是否足够
	if totalAmount < amount {
		return fmt.Errorf("insufficient balance for transfer")
	}
	// 删除原有货币
	for _, currency := range DeleteCurrencyList {
		compositeKey, err := ctx.GetStub().CreateCompositeKey("Currency", []string{oldOwner, currency.CurrencyID})
		err = ctx.GetStub().DelState(compositeKey)
		if err != nil {
			log.Fatalln("Failed to delete currency")
			return err
		}
	}
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	seconds := timestamp.GetSeconds()
	// 转账
	currencyBytes, err := json.Marshal(Currency{
		CurrencyID: "Currency" + newOwner + fmt.Sprintf("%d", seconds),
		Amount:     amount,
		Owner:      newOwner,
		CreatedAt:  fmt.Sprintf("%d", seconds),
		CreatedVia: transferReason,
		UpdatedAt:  fmt.Sprintf("%d", seconds),
		UpdatedVia: transferReason,
	})
	err = s.CreateCurrency(ctx, currencyBytes)
	if err != nil {
		return err
	}
	// 找零
	if totalAmount > amount {
		timestamp, _ := ctx.GetStub().GetTxTimestamp()
		seconds := timestamp.GetSeconds()
		currencyBytes, err := json.Marshal(Currency{
			CurrencyID: "Currency" + oldOwner + fmt.Sprintf("%d", seconds),
			Amount:     totalAmount - amount,
			Owner:      oldOwner,
			CreatedAt:  fmt.Sprintf("%d", seconds),
			CreatedVia: "Change",
			UpdatedAt:  fmt.Sprintf("%d", seconds),
			UpdatedVia: "Change",
		})
		err = s.CreateCurrency(ctx, currencyBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateContract 创建合同函数，根据业务类型，调用不同的创建合同函数
func (s *SmartContract) CreateContract(ctx contractapi.TransactionContextInterface, applicant string, businessId string, amount float32, issuer string, rate float32, businessType string, period int) error {
	switch businessType {
	case "Loan":
		return s.CreateLoan(ctx, applicant, businessId, amount, issuer, rate, period)
	case "Insurance":
		return s.CreateInsurance(ctx, applicant, businessId, amount, issuer, rate)
	default:
		return fmt.Errorf("unknown business type")
	}
}

/* Insurance 全流程
 * 保险合同结构体，用于记录保险合同的基本信息
 * CreateInsurance 创建保险合同。还未支付保险金，只是创建了保险合同。因此该函数只是创建一个“Applied”状态的保险合同。
 * ReadInsurance 读取保险合同
 * StartInsurance 保险启动函数，用于启动保险合同，支付保险金
 * InsuranceContractCheck 保险合同检查函数，检查保险是否进入赔偿状态
 */

// Insurance 保险合同结构体，用于记录保险合同的基本信息
type Insurance struct {
	BusinessID string  `json:"BusinessID"` //格式为"Insurance"+时间戳
	Amount     float32 `json:"Amount"`
	Issuer     string  `json:"Issuer"`
	State      string  `json:"State"` //"Applied","Approved","Rejected","Expired","Claimed"
	Rate       float32 `json:"Rate"`
	Applicant  string  `json:"Applicant"`
	CreatedAt  string  `json:"CreatedAt"`
	UpdatedAt  string  `json:"UpdatedAt"`
}

// CreateInsurance 创建保险合同。还未支付保险金，只是创建了保险合同。因此该函数只是创建一个“Applied”状态的保险合同。
// id 参数是保险合同的ID，应该是一个唯一的字符串，格式为"Insurance"+时间戳
func (s *SmartContract) CreateInsurance(ctx contractapi.TransactionContextInterface, applicant string, businessId string, amount float32, issuer string, rate float32) error {
	compositeKey, _ := ctx.GetStub().CreateCompositeKey("Insurance", []string{applicant, businessId})
	existing, err := s.readState(ctx, compositeKey)
	if err == nil && existing != nil {
		return fmt.Errorf("the asset %s already exists", businessId)
	}
	newTimes, _ := ctx.GetStub().GetTxTimestamp()
	seconds := newTimes.GetSeconds()
	assetJSON, err := json.Marshal(Insurance{
		BusinessID: businessId,
		Amount:     amount,
		Issuer:     issuer,
		State:      "Applied",
		Rate:       rate,
		Applicant:  applicant,
		CreatedAt:  fmt.Sprintf("%d", seconds),
		UpdatedAt:  fmt.Sprintf("%d", seconds),
	})
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("CreateInsurance", assetJSON)
	return ctx.GetStub().PutState(compositeKey, assetJSON)
}

// ReadInsurance 读取保险合同
// id 参数是保险合同的ID，是一个复合键，为Owner和BusinessID的组合
func (s *SmartContract) ReadInsurance(ctx contractapi.TransactionContextInterface, owner string, id string) (*Insurance, error) {
	compositeKey, _ := ctx.GetStub().CreateCompositeKey("Insurance", []string{owner, id})
	assetJSON, err := s.readState(ctx, compositeKey)
	if err != nil {
		return nil, err
	}

	var asset Insurance
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// StartInsurance 保险启动函数，用于启动保险合同，支付保险金
// 需要根据保险合同申请人的信用分和收入情况，判断保险是否可以启动
// 如果保险启动成功，则支付保险金，修改保险合同状态为"Approved"，并返回true
func (s *SmartContract) StartInsurance(ctx contractapi.TransactionContextInterface, applicant string, businessId string, credit float32, income float32) (bool, error) {
	compositeKey, _ := ctx.GetStub().CreateCompositeKey("Insurance", []string{applicant, businessId})
	//读取保险合同
	insurance, err := s.ReadInsurance(ctx, applicant, businessId)
	if err != nil {
		return false, err
	}
	//检查保险是否处于申请状态
	if insurance.State != "Applied" {
		return false, fmt.Errorf("the insurance contract %s is not in Applied state", businessId)
	}
	newTimes, _ := ctx.GetStub().GetTxTimestamp()
	seconds := newTimes.GetSeconds()
	//检查是否符合启动保险的条件
	if credit < 60 || income < 5000 {
		//修改保险合同状态
		insurance.State = "Rejected"
		insurance.UpdatedAt = fmt.Sprintf("%d", seconds)
		insuranceJSON, err := json.Marshal(insurance)
		if err != nil {
			return false, err
		}
		ctx.GetStub().SetEvent("StartInsurance", insuranceJSON)
		return false, ctx.GetStub().PutState(compositeKey, insuranceJSON)
	}
	//符合启动保险的条件
	//支付保险金
	err = s.TransferCurrency(ctx, insurance.Applicant, insurance.Issuer, insurance.Amount, "Insurance")
	if err != nil {
		return false, err
	}
	//修改保险合同状态
	insurance.State = "Approved"
	insurance.UpdatedAt = fmt.Sprintf("%d", seconds)
	insuranceJSON, err := json.Marshal(insurance)
	if err != nil {
		return false, fmt.Errorf("failed to marshal insurance")
	}
	ctx.GetStub().SetEvent("StartInsurance", insuranceJSON)
	return true, ctx.GetStub().PutState(compositeKey, insuranceJSON)
}

// InsuranceContractCheck 保险合同检查函数，检查保险是否进入赔偿状态
// credit 信用分，income 收入，isSudden 是否突发事件，contingencyInfo 突发事件信息
// 如果经过逻辑判断，保险需要赔偿，则立即支付赔偿金额，然后修改保险合同状态为"Claimed"，并返回true
func (s *SmartContract) InsuranceContractCheck(ctx contractapi.TransactionContextInterface, applicant string, businessId string, credit float32, income float32, isSudden bool, contingencyInfo string) (bool, error) {
	//读取保险合同
	compositeKey, _ := ctx.GetStub().CreateCompositeKey("Insurance", []string{applicant, businessId})
	insurance, err := s.ReadInsurance(ctx, applicant, businessId)
	if err != nil {
		return false, err
	}
	//检查保险是否处于申请状态
	if insurance.State != "Approved" {
		return false, fmt.Errorf("the insurance contract %s is not in Approved state", businessId)
	}
	//检查是否需要赔偿
	if credit > 60 && income < 10000 && isSudden {
		//支付赔偿
		err := s.TransferCurrency(ctx, insurance.Issuer, insurance.Applicant, insurance.Amount*(1+insurance.Rate), "Insurance")
		if err != nil {
			return false, err
		}
		newTimes, _ := ctx.GetStub().GetTxTimestamp()
		seconds := newTimes.GetSeconds()
		//未来这里可以补充对突发事件具体信息的处理逻辑

		//修改保险合同状态
		insurance.State = "Claimed"
		insurance.UpdatedAt = fmt.Sprintf("%d", seconds)
		insuranceJSON, err := json.Marshal(insurance)
		if err != nil {
			return false, err
		}
		ctx.GetStub().SetEvent("InsuranceContractCheck", insuranceJSON)
		return true, ctx.GetStub().PutState(compositeKey, insuranceJSON)
	}
	//当前不属于赔偿情况
	return false, fmt.Errorf("the insurance contract %s is not in Claimed state", businessId)
}

// ReadInsuranceListByOwner 通过owner查询保险合同列表，是一个辅助函数
func (s *SmartContract) ReadInsuranceListByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Insurance, error) {
	// 通过owner查询保险合同列表
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Insurance", []string{owner})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var insuranceList []*Insurance
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var insurance Insurance
		err = json.Unmarshal(queryResponse.Value, &insurance)
		if err != nil {
			return nil, err
		}
		if insurance.Applicant == owner {
			insuranceList = append(insuranceList, &insurance)
		}
	}
	return insuranceList, nil
}

/* Loan 全流程
 * 贷款合同结构体，用于记录贷款合同的基本信息
 * CreateLoan 创建贷款合同
 * ReadLoan 读取贷款合同
 * StartLoan 贷款启动函数，用于启动贷款合同，贷款机构向申请人支付贷款金额
 * CountLoansByOwner 通过owner查询处于”Approved“状态的贷款合同数量
 * LoanContractCheck 贷款合同检查函数，检查贷款是否进入强制还款状态
 * ReadLoanListByOwner 通过owner查询贷款合同列表，是一个辅助函数
 */

type Loan struct {
	BusinessID string  `json:"BusinessID"` //格式为"Loan"+时间戳
	Amount     float32 `json:"Amount"`
	Issuer     string  `json:"Issuer"`
	State      string  `json:"State"` //"Applied","Approved","Rejected","Expired","Claimed"
	//贷款期限
	Period int `json:"Period"`
	//贷款利率
	Rate      float32 `json:"Rate"`
	Applicant string  `json:"Applicant"`
	CreatedAt string  `json:"CreatedAt"`
	UpdatedAt string  `json:"UpdatedAt"`
}

// CreateLoan 创建贷款合同
// id 参数是贷款合同的ID，应该是一个唯一的字符串，格式为"Loan"+时间戳
func (s *SmartContract) CreateLoan(ctx contractapi.TransactionContextInterface, applicant string, businessId string, amount float32, issuer string, rate float32, period int) error {
	compositeKey, _ := ctx.GetStub().CreateCompositeKey("Loan", []string{applicant, businessId})
	existing, err := s.readState(ctx, compositeKey)
	if err == nil && existing != nil {
		return fmt.Errorf("the asset %s already exists", businessId)
	}
	newTimes, _ := ctx.GetStub().GetTxTimestamp()
	seconds := newTimes.GetSeconds()
	assetJSON, err := json.Marshal(Loan{
		BusinessID: businessId,
		Amount:     amount,
		Issuer:     issuer,
		State:      "Applied",
		Rate:       rate,
		Period:     period,
		Applicant:  applicant,
		CreatedAt:  fmt.Sprintf("%d", seconds),
		UpdatedAt:  fmt.Sprintf("%d", seconds),
	})
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("CreateLoan", assetJSON)
	return ctx.GetStub().PutState(compositeKey, assetJSON)
}

func (s *SmartContract) ReadLoan(ctx contractapi.TransactionContextInterface, owner string, id string) (*Loan, error) {
	compositeKey, _ := ctx.GetStub().CreateCompositeKey("Loan", []string{owner, id})
	assetJSON, err := s.readState(ctx, compositeKey)
	if err != nil {
		return nil, err
	}

	var asset Loan
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// CountLoansByOwner 通过owner查询处于”Approved“状态的贷款合同数量
func (s *SmartContract) CountLoansByOwner(ctx contractapi.TransactionContextInterface, owner string) (int, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Loan", []string{owner})
	if err != nil {
		return 0, err
	}
	defer resultsIterator.Close()
	var count int
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return 0, err
		}
		var loan Loan
		err = json.Unmarshal(queryResponse.Value, &loan)
		if err != nil {
			return 0, err
		}
		if loan.Applicant == owner && loan.State == "Approved" {
			count++
		}
	}
	return count, nil
}

// StartLoan 贷款启动函数，用于启动贷款合同，贷款机构向申请人支付贷款金额
// 需要根据贷款合同申请人的信用分和收入情况，判断贷款是否可以启动
// 如果贷款启动成功，则支付贷款金额，修改贷款合同状态为"Approved"，并返回true
func (s *SmartContract) StartLoan(ctx contractapi.TransactionContextInterface, applicant string, businessId string, credit float32, income float32) (bool, error) {
	compositeKey, _ := ctx.GetStub().CreateCompositeKey("Loan", []string{applicant, businessId})
	//读取贷款合同
	loan, err := s.ReadLoan(ctx, applicant, businessId)
	if err != nil {
		return false, err
	}
	//检查贷款是否处于申请状态
	if loan.State != "Applied" {
		return false, fmt.Errorf("the loan contract %s is not in Applied state", businessId)
	}
	newTimes, _ := ctx.GetStub().GetTxTimestamp()
	seconds := newTimes.GetSeconds()
	LoanCount, err := s.CountLoansByOwner(ctx, applicant)
	//检查是否符合启动贷款的条件
	if credit < 60 || income < 5000 || loan.Amount > 10000 || LoanCount > 3 {
		//修改贷款合同状态
		loan.State = "Rejected"
		loan.UpdatedAt = fmt.Sprintf("%d", seconds)
		loanJSON, err := json.Marshal(loan)
		if err != nil {
			return false, err
		}
		ctx.GetStub().SetEvent("StartLoan", loanJSON)
		return false, ctx.GetStub().PutState(compositeKey, loanJSON)
	}
	//符合启动贷款的条件
	//支付贷款金额
	err = s.TransferCurrency(ctx, loan.Issuer, loan.Applicant, loan.Amount, "Loan")
	if err != nil {
		return false, err
	}
	//修改贷款合同状态
	loan.State = "Approved"
	loan.UpdatedAt = fmt.Sprintf("%d", seconds)
	loanJSON, err := json.Marshal(loan)
	if err != nil {
		return false, fmt.Errorf("failed to marshal loan")
	}
	ctx.GetStub().SetEvent("StartLoan", loanJSON)
	return true, ctx.GetStub().PutState(compositeKey, loanJSON)
}

// LoanContractCheck 贷款合同检查函数，检查贷款是否进入强制还款状态
// credit 信用分，income 收入，isOverdue 是否逾期
// 如果经过逻辑判断，贷款需要强制还款，则立即支付剩余贷款金额，然后修改贷款合同状态为"Claimed"，并返回true
func (s *SmartContract) LoanContractCheck(ctx contractapi.TransactionContextInterface, applicant string, businessId string, credit float32, income float32, currentTime string) (bool, error) {
	var isOverdue = false
	//读取贷款合同
	compositeKey, _ := ctx.GetStub().CreateCompositeKey("Loan", []string{applicant, businessId})
	loan, err := s.ReadLoan(ctx, applicant, businessId)
	if err != nil {
		return false, err
	}
	//检查贷款是否处于申请状态
	if loan.State != "Approved" {
		return false, fmt.Errorf("the loan contract %s is not in Approved state", businessId)
	}
	//判断是否逾期
	currentTimestamp, _ := strconv.Atoi(currentTime)
	loanTimestamp, _ := strconv.Atoi(loan.CreatedAt)
	if currentTimestamp-loanTimestamp > loan.Period*24*60*60 {
		isOverdue = true
	}
	//检查是否需要强制还款
	if credit > 60 || income < 5000 || isOverdue {
		//支付剩余贷款
		err := s.TransferCurrency(ctx, loan.Applicant, loan.Issuer, loan.Amount*(1+loan.Rate), "Loan")
		if err != nil {
			return false, err
		}
		newTimes, _ := ctx.GetStub().GetTxTimestamp()
		seconds := newTimes.GetSeconds()
		//修改贷款合同状态
		loan.State = "Claimed"
		loan.UpdatedAt = fmt.Sprintf("%d", seconds)
		loanJSON, err := json.Marshal(loan)
		if err != nil {
			return false, err
		}
		ctx.GetStub().SetEvent("LoanContractCheck", loanJSON)
		return true, ctx.GetStub().PutState(compositeKey, loanJSON)
	}
	//当前不属于强制还款情况
	return false, fmt.Errorf("the loan contract %s is not in Claimed state", businessId)
}

// ReadLoanListByOwner 根据owner读取所有的贷款合同
func (s *SmartContract) ReadLoanListByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Loan, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Loan", []string{owner})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var loans []*Loan
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var loan Loan
		err = json.Unmarshal(queryResponse.Value, &loan)
		if err != nil {
			return nil, err
		}
		if loan.Applicant == owner {
			loans = append(loans, &loan)
		}
	}
	return loans, nil
}

/* 其他辅助链码函数
 * readState 通用的资产结构（未解析）读取函数
 */

// 通用的资产结构（未解析）读取函数
func (s *SmartContract) readState(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	return assetJSON, nil
}

// Asset describes basic details of what makes up a simple asset【asset部分是原有的，用于测试/参考，并不作为实险贷链码的一部分】
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	AppraisedValue int    `json:"AppraisedValue"`
	Color          string `json:"Color"`
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	Size           int    `json:"Size"`
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	existing, err := s.readState(ctx, id)
	if err == nil && existing != nil {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("CreateAsset", assetJSON)
	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := s.readState(ctx, id)
	if err != nil {
		return nil, err
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	_, err := s.readState(ctx, id)
	if err != nil {
		return err
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("UpdateAsset", assetJSON)
	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes a given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	assetJSON, err := s.readState(ctx, id)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("DeleteAsset", assetJSON)
	return ctx.GetStub().DelState(id)
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	ctx.GetStub().SetEvent("TransferAsset", assetJSON)
	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
