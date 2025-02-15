// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

//nolint:gocritic
package zentangle

import "github.com/iotaledger/wasp/packages/vm/wasmlib/go/wasmlib"

type zentangleEvents struct {
}

func (e zentangleEvents) GameEnded() {
	wasmlib.NewEventEncoder("zentangle.gameEnded").
		Emit()
}

func (e zentangleEvents) GameStarted(description string, numberOfImages int32, reward int64, tagsRequiredPerImage int32) {
	wasmlib.NewEventEncoder("zentangle.gameStarted").
		String(description).
		Int32(numberOfImages).
		Int64(reward).
		Int32(tagsRequiredPerImage).
		Emit()
}

func (e zentangleEvents) Imagetagged(address string, imageId int32, playsPerImage int32) {
	wasmlib.NewEventEncoder("zentangle.imagetagged").
		String(address).
		Int32(imageId).
		Int32(playsPerImage).
		Emit()
}

func (e zentangleEvents) PlayRequested(address string, amount int64, imageId int32) {
	wasmlib.NewEventEncoder("zentangle.playRequested").
		String(address).
		Int64(amount).
		Int32(imageId).
		Emit()
}
