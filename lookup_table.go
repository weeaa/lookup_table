package lookup_table

import (
	"bytes"
	"encoding/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

var PROGRAM_ID = solana.MustPublicKeyFromBase58("AddressLookupTab1e1111111111111111111111111")

func Create(recentSlot uint64, authority, payer solana.PublicKey) (*solana.GenericInstruction, solana.PublicKey, error) {
	recentSlotBytes := toBufferLittleEndian(recentSlot)

	address, seed, err := solana.FindProgramAddress(
		[][]byte{
			authority.Bytes(),
			recentSlotBytes,
		},
		PROGRAM_ID,
	)
	if err != nil {
		return nil, solana.PublicKey{}, err
	}

	data := []byte{0}
	data = append(data, recentSlotBytes...)
	data = append(data, seed)

	keys := []*solana.AccountMeta{
		{PublicKey: address, IsSigner: false, IsWritable: true},
		{PublicKey: authority, IsSigner: true, IsWritable: false},
		{PublicKey: payer, IsSigner: true, IsWritable: true},
		{PublicKey: system.ProgramID, IsSigner: false, IsWritable: false},
	}

	return solana.NewInstruction(
		PROGRAM_ID,
		keys,
		data,
	), address, nil
}

func Extend(tableAddress, authority solana.PublicKey, payer *solana.PublicKey, addresses []solana.PublicKey) (*solana.GenericInstruction, error) {
	data := []byte{2}
	data = append(data, toBufferLittleEndian(uint64(len(addresses)))...)
	for _, address := range addresses {
		data = append(data, address.Bytes()...)
	}

	keys := []*solana.AccountMeta{
		{PublicKey: tableAddress, IsSigner: false, IsWritable: true},
		{PublicKey: authority, IsSigner: true, IsWritable: false},
	}
	if payer != nil {
		keys = append(keys, []*solana.AccountMeta{
			{PublicKey: *payer, IsSigner: true, IsWritable: true},
			{PublicKey: system.ProgramID, IsSigner: false, IsWritable: false},
		}...)
	}

	return solana.NewInstruction(
		PROGRAM_ID,
		keys,
		data,
	), nil
}

func Close(tableAddress, authority, recipient solana.PublicKey) (*solana.GenericInstruction, error) {
	data := []byte{4}

	keys := []*solana.AccountMeta{
		{PublicKey: tableAddress, IsSigner: false, IsWritable: true},
		{PublicKey: authority, IsSigner: true, IsWritable: false},
		{PublicKey: recipient, IsSigner: false, IsWritable: true},
	}

	return solana.NewInstruction(
		PROGRAM_ID,
		keys,
		data,
	), nil
}

func Freeze(tableAddress, authority solana.PublicKey) (*solana.GenericInstruction, error) {
	data := []byte{1}

	keys := []*solana.AccountMeta{
		{PublicKey: tableAddress, IsSigner: false, IsWritable: true},
		{PublicKey: authority, IsSigner: true, IsWritable: false},
	}

	return solana.NewInstruction(
		PROGRAM_ID,
		keys,
		data,
	), nil
}

func Deactivate(tableAddress, authority solana.PublicKey) (*solana.GenericInstruction, error) {
	data := []byte{3}

	keys := []*solana.AccountMeta{
		{PublicKey: tableAddress, IsSigner: false, IsWritable: true},
		{PublicKey: authority, IsSigner: true, IsWritable: false},
	}

	return solana.NewInstruction(
		PROGRAM_ID,
		keys,
		data,
	), nil
}

func toBufferLittleEndian(value uint64) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	return buf.Bytes()
}
