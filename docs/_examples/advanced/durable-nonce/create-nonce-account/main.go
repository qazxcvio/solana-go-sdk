package main

import (
	"context"
	"fmt"
	"log"

	"github.com/qazxcvio/solana-go-sdk/client"
	"github.com/qazxcvio/solana-go-sdk/common"
	"github.com/qazxcvio/solana-go-sdk/program/system"
	"github.com/qazxcvio/solana-go-sdk/rpc"
	"github.com/qazxcvio/solana-go-sdk/types"
)

// FUarP2p5EnxD66vVDL4PWRoWMzA56ZVHG24hpEDFShEz
var feePayer, _ = types.AccountFromBase58("4TMFNY9ntAn3CHzguSAvDNLPRoQTaK3sWbQQXdDXaE6KWRBLufGL6PJdsD2koiEe3gGmMdRK3aAw7sikGNksHJrN")

// 9aE476sH92Vz7DMPyq5WLPkrKWivxeuTKEFKd2sZZcde
var alice, _ = types.AccountFromBase58("4voSPg3tYuWbKzimpQK9EbXHmuyy5fUrtXvpLDMLkmY6TRncaTHAKGD8jUg3maB5Jbrd9CkQg4qjJMyN6sQvnEF2")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	// create a new account
	nonceAccount := types.NewAccount()
	fmt.Println("nonce account:", nonceAccount.PublicKey)

	// get minimum balance
	nonceAccountMinimumBalance, err := c.GetMinimumBalanceForRentExemption(context.Background(), system.NonceAccountSize)
	if err != nil {
		log.Fatalf("failed to get minimum balance for nonce account, err: %v", err)
	}

	// recent blockhash
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	// create a tx
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer, nonceAccount},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				system.CreateAccount(system.CreateAccountParam{
					From:     feePayer.PublicKey,
					New:      nonceAccount.PublicKey,
					Owner:    common.SystemProgramID,
					Lamports: nonceAccountMinimumBalance,
					Space:    system.NonceAccountSize,
				}),
				system.InitializeNonceAccount(system.InitializeNonceAccountParam{
					Nonce: nonceAccount.PublicKey,
					Auth:  alice.PublicKey,
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
	}

	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	fmt.Println("txhash", sig)
}
