package internal

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"strings"
)

const (
	StakeCommand    = "stake"
	WithdrawCommand = "withdraw"
	ClaimCommand    = "claim"
)

type StakingCommand struct {
	Amount     *big.Int
	PrivateKey string
	Address    string
	Command    string
	ConfigFile string
	TxHash     string
}

type StakingClient struct {
	client  *ethclient.Client
	config  *Config
	command *StakingCommand
}

func (sc *StakingCommand) Run() error {
	stakingClient, err := NewStakingClient(sc)
	if err != nil {
		log.Println("Error creating staking client:", err)
	}

	switch sc.Command {
	case StakeCommand:
		return stakingClient.Stake()
	case WithdrawCommand:
		return stakingClient.Withdraw()
	case ClaimCommand:
		return stakingClient.Claim()
	}
	return nil
}

func NewStakingClient(command *StakingCommand) (*StakingClient, error) {
	config, err := LoadConfig(command.ConfigFile)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(config.AppConfig.RPCURL)
	if err != nil {
		return nil, err
	}

	return &StakingClient{
		client:  client,
		config:  config,
		command: command,
	}, nil
}

func (scl *StakingClient) Stake() error {
	data, err := scl.GetData("stake", scl.config.AppConfig.StakeABI)
	if err != nil {
		return err
	}

	return scl.ExecuteTx(scl.command.Amount, data)
}

func (scl *StakingClient) Withdraw() error {
	data, err := scl.GetData("withdrawAll", scl.config.AppConfig.WithdrawABI)
	if err != nil {
		return err
	}

	return scl.ExecuteTx(big.NewInt(0), data)
}

func (scl *StakingClient) Claim() error {
	data, err := scl.GetData("claimReward", scl.config.AppConfig.ClaimABI)
	if err != nil {
		return err
	}

	return scl.ExecuteTx(big.NewInt(0), data)
}

func (scl *StakingClient) GetData(funcName string, funcAbi string) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(funcAbi))
	if err != nil {
		log.Fatalf("ABI parsing error: %v", err)
		return nil, err
	}

	data, err := parsedABI.Pack(funcName)
	if err != nil {
		log.Fatalf("Data packing error: %v", err)
		return nil, err
	}

	return data, nil
}

func (scl *StakingClient) ExecuteTx(value *big.Int, data []byte) error {
	cleanedKey := strings.TrimPrefix(scl.command.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(cleanedKey)
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
		return err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	nonce, err := scl.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Error getting nonce: %v", err)
		return err
	}

	toAddress := common.HexToAddress(scl.config.AppConfig.ContractAddress)

	tipCap, err := scl.client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal("Error getting tipCap:", err)
		return err
	}

	feeCap, err := scl.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Error getting feeCap:", err)
		return err
	}

	gasLimit, err := scl.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: value,
		Data:  data,
	})
	if err != nil {
		log.Fatalf("Ошибка расчета газа: %v", err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &toAddress,
		Value:     value,
		Gas:       gasLimit,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Data:      data,
	})

	chainID, ok := new(big.Int).SetString(scl.config.AppConfig.ChainID, 10)

	if !ok {
		log.Fatal("Chain Id Incorrect:", err)
		return err
	}

	signer := types.NewLondonSigner(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatalf("Transaction signature error: %v", err)
		return err
	}
	signedTx, err = types.SignTx(tx, signer, privateKey)

	err = scl.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Error sending transaction: %v", err)
		return err
	}

	scl.command.TxHash = signedTx.Hash().Hex()
	scl.logTransactionHash()

	return nil
}

func (scl *StakingClient) logTransactionHash() {
	file, err := logFile()
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(scl.command.TxHash + "\n")
	if err != nil {
		log.Fatalf("Error writing to log file: %v", err)
	}
}

func logFile() (*os.File, error) {
	return os.OpenFile("transactions.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
