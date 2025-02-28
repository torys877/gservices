# Staking Tool CLI

The **Staking Tool CLI** is a command-line interface (CLI) application designed to interact with a staking smart contract on the Sepolia blockchain. 
It allows users to stake ETH, withdraw staked ETH and rewards, and claim rewards directly from the terminal. 
The tool is built using Go and leverages the `cobra` library for CLI functionality and the `go-ethereum` library for blockchain interactions.

## Staking address on Sepolia

The staking address to interact with on Sepolia is: `0x0ec6d8a992B599A6413bd8f4241dE22a9070EFf6`

---

## Features

- **Stake ETH**: Stake a specified amount of ETH into the staking contract.
- **Withdraw ETH**: Withdraw all staked ETH and accumulated rewards from the contract.
- **Claim Rewards**: Claim rewards without withdrawing the staked ETH.
- **Configurable**: Use a YAML configuration file to specify RPC URL, chain ID, contract address, and ABIs.
- **Transaction Logging**: Logs transaction hashes to a file for easy tracking.

---

## Prerequisites

Before using the Staking Tool CLI, ensure you have the following:

1. **Go** installed (version 1.23).
2. A **private key** for the Ethereum account you want to use.
3. Access to an Ethereum node or a service like Infura/Alchemy for the RPC URL.
4. A YAML configuration file with the required settings (see [Configuration](#configuration)).

---

## Preparation
You need to get some Ethers from any public Sepolia faucet, eg:

```
https://cloud.google.com/application/web3/faucet/ethereum/sepolia
```

## Installation

1. Clone the repository or download the source code.

```bash
git clone <repository-url>
cd staketool
``` 

2. Install dependencies.

```bash 
go mod tidy
```

2. Build the project.

```bash
go build -o staketool
```

## Usage

1. Stake Ethers

```bash
./staketool stake --config path/to/config.yaml --amount <amount> --private_key <your_private_key>
```

or using Go directly:

```bash
go run .\cmd\main.go stake --amount <amount> --private_key <your_private_key> --config .\config\config.yaml
```

2. Claim rewards

```bash
./staketool claim --config path/to/config.yaml --private_key <your_private_key>
```

or using Go directly:

```bash
go run .\cmd\main.go claim --config path/to/config.yaml --private_key <your_private_key>
```

3. Withdraw all staked ETH and rewards from the contract

```bash
./staketool withdraw --config path/to/config.yaml --private_key <your_private_key>
```

or using Go directly:

```bash
go run .\cmd\main.go withdraw --config path/to/config.yaml --private_key <your_private_key>
```

## Configuration
The CLI requires a YAML configuration file (config.yaml) with the following structure:

```yaml
config:
  rpc_url: "<RPC_URL>" # SET your sepolia RPC URL
  chain_id: "11155111" # spolia chain_id
  contract_address: "0x0ec6d8a992B599A6413bd8f4241dE22a9070EFf6" # sepolia staking smart contract
  stake_abi: '[{
      "inputs": [],
      "name": "stake",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    }]'
  withdraw_abi: '[{
      "inputs": [],
      "name": "withdrawAll",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }]'
  claim_abi: '[{
      "inputs": [],
      "name": "claimReward",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }]'

```

> You need to set your own Sepolia `rpc_url`

## How It Works
1. **Command Parsing**: The CLI uses the cobra library to parse commands and flags.
2. **Configuration Loading**: The YAML configuration file is loaded to retrieve RPC URL, chain ID, contract address, and ABIs.
3. **Transaction Execution**:
   * The private key is used to sign transactions.
   * Gas fees and limits are estimated dynamically.
   * Transactions are sent to the Ethereum network using the go-ethereum library.
4. **Logging**: Transaction hashes are logged to a file (transactions.log) for future reference.

## Logs
All transaction hashes are logged to transactions.log in the project directory. Example:

```
0x123...abc
0x456...def
```
