// Copyright (C) 2025 Bitfly GmbH
//
// This file is part of Beaconchain Dashboard.
//
// Beaconchain Dashboard is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Beaconchain Dashboard is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Beaconchain Dashboard.  If not, see <https://www.gnu.org/licenses/>.

package weth

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestWethABI verifies that the WETH ABI can be parsed correctly
func TestWethABI(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(WethABI))
	if err != nil {
		t.Fatalf("Failed to parse WETH ABI: %v", err)
	}

	// Verify that key methods exist
	methods := []string{"transfer", "transferFrom", "deposit", "withdraw", "balanceOf", "approve"}
	for _, method := range methods {
		if _, ok := parsedABI.Methods[method]; !ok {
			t.Errorf("Method %s not found in WETH ABI", method)
		}
	}

	// Verify that key events exist
	events := []string{"Transfer", "Approval", "Deposit", "Withdrawal"}
	for _, event := range events {
		if _, ok := parsedABI.Events[event]; !ok {
			t.Errorf("Event %s not found in WETH ABI", event)
		}
	}
}

// TestTransferEventSignature verifies the Transfer event signature
func TestTransferEventSignature(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(WethABI))
	if err != nil {
		t.Fatalf("Failed to parse WETH ABI: %v", err)
	}

	transferEvent, ok := parsedABI.Events["Transfer"]
	if !ok {
		t.Fatal("Transfer event not found in ABI")
	}

	// The Transfer event signature should be: Transfer(address,address,uint256)
	expectedSig := "Transfer(address,address,uint256)"
	actualSig := transferEvent.Sig

	if actualSig != expectedSig {
		t.Errorf("Transfer event signature mismatch: got %s, want %s", actualSig, expectedSig)
	}

	// Verify the event ID (topic0)
	expectedID := crypto.Keccak256Hash([]byte(expectedSig))
	actualID := transferEvent.ID

	if actualID != expectedID {
		t.Errorf("Transfer event ID mismatch: got %s, want %s", actualID.Hex(), expectedID.Hex())
	}
}

// TestTransferMethodEncoding verifies the transfer method can be encoded
func TestTransferMethodEncoding(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(WethABI))
	if err != nil {
		t.Fatalf("Failed to parse WETH ABI: %v", err)
	}

	transferMethod, ok := parsedABI.Methods["transfer"]
	if !ok {
		t.Fatal("transfer method not found in ABI")
	}

	// Test encoding a transfer call
	recipient := common.HexToAddress("0x1234567890123456789012345678901234567890")
	amount := big.NewInt(1000000000000000000) // 1 WETH

	data, err := transferMethod.Inputs.Pack(recipient, amount)
	if err != nil {
		t.Fatalf("Failed to pack transfer arguments: %v", err)
	}

	if len(data) == 0 {
		t.Error("Packed transfer data is empty")
	}

	// Verify we can unpack the data
	unpacked, err := transferMethod.Inputs.Unpack(data)
	if err != nil {
		t.Fatalf("Failed to unpack transfer data: %v", err)
	}

	if len(unpacked) != 2 {
		t.Errorf("Expected 2 unpacked arguments, got %d", len(unpacked))
	}
}

// TestTransferFromMethodEncoding verifies the transferFrom method can be encoded
func TestTransferFromMethodEncoding(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(WethABI))
	if err != nil {
		t.Fatalf("Failed to parse WETH ABI: %v", err)
	}

	transferFromMethod, ok := parsedABI.Methods["transferFrom"]
	if !ok {
		t.Fatal("transferFrom method not found in ABI")
	}

	// Test encoding a transferFrom call
	sender := common.HexToAddress("0x1111111111111111111111111111111111111111")
	recipient := common.HexToAddress("0x2222222222222222222222222222222222222222")
	amount := big.NewInt(500000000000000000) // 0.5 WETH

	data, err := transferFromMethod.Inputs.Pack(sender, recipient, amount)
	if err != nil {
		t.Fatalf("Failed to pack transferFrom arguments: %v", err)
	}

	if len(data) == 0 {
		t.Error("Packed transferFrom data is empty")
	}

	// Verify we can unpack the data
	unpacked, err := transferFromMethod.Inputs.Unpack(data)
	if err != nil {
		t.Fatalf("Failed to unpack transferFrom data: %v", err)
	}

	if len(unpacked) != 3 {
		t.Errorf("Expected 3 unpacked arguments, got %d", len(unpacked))
	}
}

// TestTransferEventDecoding verifies that Transfer events can be decoded
func TestTransferEventDecoding(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(WethABI))
	if err != nil {
		t.Fatalf("Failed to parse WETH ABI: %v", err)
	}

	transferEvent, ok := parsedABI.Events["Transfer"]
	if !ok {
		t.Fatal("Transfer event not found in ABI")
	}

	// Create test event data
	amount := big.NewInt(1000000000000000000) // 1 WETH

	// Pack the non-indexed data (amount)
	data, err := transferEvent.Inputs.NonIndexed().Pack(amount)
	if err != nil {
		t.Fatalf("Failed to pack event data: %v", err)
	}

	// Verify we can unpack the data
	var decoded struct {
		Wad *big.Int
	}

	err = parsedABI.UnpackIntoInterface(&decoded, "Transfer", data)
	if err != nil {
		t.Fatalf("Failed to unpack Transfer event: %v", err)
	}

	if decoded.Wad.Cmp(amount) != 0 {
		t.Errorf("Decoded amount mismatch: got %s, want %s", decoded.Wad.String(), amount.String())
	}
}

// TestDepositAndWithdrawMethods verifies deposit and withdraw methods exist
func TestDepositAndWithdrawMethods(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(WethABI))
	if err != nil {
		t.Fatalf("Failed to parse WETH ABI: %v", err)
	}

	// Test deposit method
	depositMethod, ok := parsedABI.Methods["deposit"]
	if !ok {
		t.Fatal("deposit method not found in ABI")
	}

	if !depositMethod.IsPayable() {
		t.Error("deposit method should be payable")
	}

	// Test withdraw method
	withdrawMethod, ok := parsedABI.Methods["withdraw"]
	if !ok {
		t.Fatal("withdraw method not found in ABI")
	}

	// Verify withdraw takes a uint256 parameter
	if len(withdrawMethod.Inputs) != 1 {
		t.Errorf("withdraw method should have 1 input, got %d", len(withdrawMethod.Inputs))
	}
}

// TestWethBindingsExist verifies that the generated bindings are available
func TestWethBindingsExist(t *testing.T) {
	// This test verifies that the auto-generated types exist
	// by attempting to reference them

	// These types should be available from bindings.go
	var _ *Weth
	var _ *WethCaller
	var _ *WethTransactor
	var _ *WethFilterer
	var _ *WethSession
	var _ *WethCallerSession
	var _ *WethTransactorSession
}
