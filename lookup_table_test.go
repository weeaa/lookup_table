package lookup_table

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	rpcClient  *rpc.Client
	privateKey solana.PrivateKey
)

func init() {
	godotenv.Load()
	key, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		panic("could not locate PRIVATE_KEY in .env")
	}
	rpcClient = rpc.New(rpc.DevNet_RPC)
	privateKey = solana.MustPrivateKeyFromBase58(key)
}

func TestCreateAndExtendTable(t *testing.T) {
	ctx := context.Background()

	slot, err := rpcClient.GetSlot(ctx, rpc.CommitmentFinalized)
	require.NoError(t, err)

	createInstruction, tableAddress, err := Create(
		slot,
		privateKey.PublicKey(),
		privateKey.PublicKey(),
	)
	require.NoError(t, err)

	recent, err := rpcClient.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	require.NoError(t, err)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{createInstruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	require.NoError(t, err)

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if privateKey.PublicKey().Equals(key) {
			return &privateKey
		}
		return nil
	})
	require.NoError(t, err)

	sig, err := rpcClient.SendTransaction(
		ctx,
		tx,
	)
	require.NoError(t, err)

	time.Sleep(time.Second * 2)
	fmt.Printf("create table sig: %s", sig.String())

	addresses := []solana.PublicKey{
		solana.NewWallet().PublicKey(),
		solana.NewWallet().PublicKey(),
	}

	pkPtr := privateKey.PublicKey()
	extendInstruction, err := Extend(
		tableAddress,
		privateKey.PublicKey(),
		&pkPtr,
		addresses,
	)
	require.NoError(t, err)

	recent, err = rpcClient.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	require.NoError(t, err)

	extendTx, err := solana.NewTransaction(
		[]solana.Instruction{extendInstruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	require.NoError(t, err)

	_, err = extendTx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if privateKey.PublicKey().Equals(key) {
			return &privateKey
		}
		return nil
	})
	require.NoError(t, err)

	sig, err = rpcClient.SendTransaction(
		ctx,
		extendTx,
	)
	require.NoError(t, err)
	fmt.Printf("extend table sig: %s", sig.String())

	accountInfo, err := rpcClient.GetAccountInfo(ctx, tableAddress)
	require.NoError(t, err)
	require.NotNil(t, accountInfo)
}

func TestExtendTable(t *testing.T) {
	ctx := context.Background()
	tableAddress := solana.MustPublicKeyFromBase58("") // rm todo add in env

	addresses := []solana.PublicKey{
		solana.NewWallet().PublicKey(),
		solana.NewWallet().PublicKey(),
	}

	pkPtr := privateKey.PublicKey()
	extendInstruction, err := Extend(
		tableAddress,
		privateKey.PublicKey(),
		&pkPtr,
		addresses,
	)
	require.NoError(t, err)

	recent, err := rpcClient.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	require.NoError(t, err)

	extendTx, err := solana.NewTransaction(
		[]solana.Instruction{extendInstruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	require.NoError(t, err)

	_, err = extendTx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if privateKey.PublicKey().Equals(key) {
			return &privateKey
		}
		return nil
	})
	require.NoError(t, err)

	sig, err := rpcClient.SendTransaction(
		ctx,
		extendTx,
	)
	require.NoError(t, err)
	fmt.Printf("extend table sig: %s", sig.String())

	accountInfo, err := rpcClient.GetAccountInfo(ctx, tableAddress)
	require.NoError(t, err)
	require.NotNil(t, accountInfo)
}
