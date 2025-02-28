package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"math/big"
	"staketool/internal"
)

var (
	privateKey string
	amount     string
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "staketool",
	Short: "CLI for interacting with the staking contract",
}

var stakeCmd = &cobra.Command{
	Use:   "stake",
	Short: "Stake ETH",
	Run:   commandHandler,
}

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Withdraw staked ETH and rewards",
	Run:   commandHandler,
}

var claimCmd = &cobra.Command{
	Use:   "claim",
	Short: "Claim rewards",
	Run:   commandHandler,
}

func commandHandler(cmd *cobra.Command, args []string) {
	if privateKey == "" {
		log.Fatal("Private key is required")
	}

	stakeCommand := internal.StakingCommand{
		ConfigFile: configFile,
		PrivateKey: privateKey,
		Command:    cmd.Name(),
	}

	if cmd.Name() == "stake" {
		amountBigInt, ok := new(big.Int).SetString(amount, 10)
		if !ok {
			log.Fatal("Invalid amount")
		}

		stakeCommand.Amount = amountBigInt
	}

	err := stakeCommand.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s command executed successfuly\n", cmd.Name())
	fmt.Printf("TxHash: %s", stakeCommand.TxHash)
}

func init() {
	stakeCmd.Flags().StringVarP(&privateKey, "private_key", "k", "", "private key for the transaction")
	stakeCmd.Flags().StringVarP(&amount, "amount", "a", "", "Amount to stake")
	stakeCmd.MarkFlagRequired("private_key")
	stakeCmd.MarkFlagRequired("amount")
	stakeCmd.MarkFlagRequired("config")

	withdrawCmd.Flags().StringVarP(&privateKey, "private_key", "k", "", "private key for the transaction")
	withdrawCmd.MarkFlagRequired("private_key")

	claimCmd.Flags().StringVarP(&privateKey, "private_key", "k", "", "private key for the transaction")
	claimCmd.MarkFlagRequired("private_key")

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "path to config file")
	rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.AddCommand(stakeCmd)
	rootCmd.AddCommand(withdrawCmd)
	rootCmd.AddCommand(claimCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
