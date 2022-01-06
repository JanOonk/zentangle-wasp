// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmsolo

import (
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/colored"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/wasmhost"
	"github.com/iotaledger/wasp/packages/vm/wasmproc"
)

type SoloScContext struct {
	wasmproc.ScContext
	ctx *SoloContext
}

func NewSoloScContext(ctx *SoloContext) *SoloScContext {
	return &SoloScContext{ScContext: *wasmproc.NewScContext(nil, &ctx.wc.KvStoreHost), ctx: ctx}
}

func (o *SoloScContext) Exists(keyID, typeID int32) bool {
	return o.GetTypeID(keyID) > 0
}

func (o *SoloScContext) GetBytes(keyID, typeID int32) []byte {
	switch keyID {
	case wasmhost.KeyChainID:
		return o.ctx.Chain.ChainID.Bytes()
	default:
		o.InvalidKey(keyID)
		return nil
	}
}

func (o *SoloScContext) GetObjectID(keyID, typeID int32) int32 {
	host := o.ctx.wc
	return wasmproc.GetMapObjectID(o, keyID, typeID, wasmproc.ObjFactories{
		// wasmhost.KeyBalances:  func() wasmproc.WaspObject { return wasmproc.NewScBalances(o.vm, keyID) },
		wasmhost.KeyExports: func() wasmproc.WaspObject { return wasmproc.NewScExports(host) },
		// wasmhost.KeyIncoming:  func() wasmproc.WaspObject { return wasmproc.NewScBalances(o.vm, keyID) },
		wasmhost.KeyMaps: func() wasmproc.WaspObject { return wasmproc.NewScMaps(&host.KvStoreHost) },
		// wasmhost.KeyMinted:    func() wasmproc.WaspObject { return wasmproc.NewScBalances(o.vm, keyID) },
		// wasmhost.KeyParams:    func() wasmproc.WaspObject { return wasmproc.NewScDict(o.host, o.vm.params()) },
		wasmhost.KeyResults: func() wasmproc.WaspObject { return wasmproc.NewScDict(&host.KvStoreHost, dict.New()) },
		wasmhost.KeyReturn:  func() wasmproc.WaspObject { return wasmproc.NewScDict(&host.KvStoreHost, dict.New()) },
		// wasmhost.KeyState:     func() wasmproc.WaspObject { return wasmproc.NewScDict(o.host, o.vm.state()) },
		// wasmhost.KeyTransfers: func() wasmproc.WaspObject { return wasmproc.NewScTransfers(o.vm) },
		wasmhost.KeyUtility: func() wasmproc.WaspObject { return wasmproc.NewScUtility(nil, o.ctx) },
	})
}

func (o *SoloScContext) SetBytes(keyID, typeID int32, bytes []byte) {
	switch keyID {
	case wasmhost.KeyCall:
		o.processCall(bytes)
	case wasmhost.KeyPost:
		o.processPost(bytes)
	case wasmhost.KeyLog:
		o.ctx.Chain.Log.Infof(string(bytes))
	case wasmhost.KeyTrace:
		o.ctx.Chain.Log.Debugf(string(bytes))
	default:
		o.ScContext.SetBytes(keyID, typeID, bytes)
	}
}

func (o *SoloScContext) processCall(bytes []byte) {
	decode := wasmproc.NewBytesDecoder(bytes)
	contract, err := iscp.HnameFromBytes(decode.Bytes())
	if err != nil {
		o.Panicf(err.Error())
	}
	function, err := iscp.HnameFromBytes(decode.Bytes())
	if err != nil {
		o.Panicf(err.Error())
	}
	paramsID := decode.Int32()
	transferID := decode.Int32()
	if transferID != 0 {
		o.postSync(contract, function, paramsID, transferID, 0)
		return
	}

	ctx := o.ctx
	funcName := ctx.wc.FunctionFromCode(uint32(function))
	if funcName == "" {
		o.Panicf("unknown function")
	}
	o.Tracef("CALL %s.%s", ctx.scName, funcName)
	params := o.getParams(paramsID)
	_ = wasmhost.Connect(ctx.wasmHostOld)
	res, err := ctx.Chain.CallView(ctx.scName, funcName, params)
	_ = wasmhost.Connect(ctx.wc)
	ctx.Err = err
	if err != nil {
		// o.Panicf("failed to invoke call: " + err.Error())
		return
	}
	returnID := o.GetObjectID(wasmhost.KeyReturn, wasmhost.OBJTYPE_MAP)
	ctx.wc.FindObject(returnID).(*wasmproc.ScDict).SetKvStore(res)
}

func (o *SoloScContext) processPost(bytes []byte) {
	decode := wasmproc.NewBytesDecoder(bytes)
	chainID, err := iscp.ChainIDFromBytes(decode.Bytes())
	if err != nil {
		o.Panicf(err.Error())
	}
	if !chainID.Equals(o.ctx.Chain.ChainID) {
		o.Panicf("invalid chainID")
	}
	contract, err := iscp.HnameFromBytes(decode.Bytes())
	if err != nil {
		o.Panicf(err.Error())
	}
	function, err := iscp.HnameFromBytes(decode.Bytes())
	if err != nil {
		o.Panicf(err.Error())
	}
	paramsID := decode.Int32()
	transferID := decode.Int32()
	delay := decode.Int32()
	o.postSync(contract, function, paramsID, transferID, delay)
	//metadata := &iscp.SendMetadata{
	//	TargetContract: contract,
	//	EntryPoint:     function,
	//	Args:           params,
	//}
	//delay := decode.Int32()
	//if delay == 0 {
	//	if !o.vm.ctx.Send(chainID.AsAddress(), transfer, metadata) {
	//		o.Panicf("failed to send to %s", chainID.AsAddress().String())
	//	}
	//	return
	//}
	//
	//if delay < -1 {
	//	o.Panicf("invalid delay: %d", delay)
	//}
	//
	//timeLock := time.Unix(0, o.vm.ctx.GetTimestamp())
	//timeLock = timeLock.Add(time.Duration(delay) * time.Second)
	//options := iscp.SendOptions{
	//	TimeLock: uint32(timeLock.Unix()),
	//}
	//if !o.vm.ctx.Send(chainID.AsAddress(), transfer, metadata, options) {
	//	o.Panicf("failed to send to %s", chainID.AsAddress().String())
	//}
}

func (o *SoloScContext) getParams(paramsID int32) dict.Dict {
	if paramsID == 0 {
		return dict.New()
	}
	params := o.ctx.wc.FindObject(paramsID).(*wasmproc.ScDict).KvStore().(dict.Dict)
	params.MustIterate("", func(key kv.Key, value []byte) bool {
		o.Tracef("  PARAM '%s'", key)
		return true
	})
	return params
}

func (o *SoloScContext) getTransfer(transferID int32) colored.Balances {
	if transferID == 0 {
		return colored.NewBalances()
	}
	transfer := colored.NewBalances()
	transferDict := o.ctx.wc.FindObject(transferID).(*wasmproc.ScDict).KvStore()
	transferDict.MustIterate("", func(key kv.Key, value []byte) bool {
		color, err := codec.DecodeColor([]byte(key))
		if err != nil {
			o.Panicf(err.Error())
		}
		amount, err := codec.DecodeUint64(value)
		if err != nil {
			o.Panicf(err.Error())
		}
		o.Tracef("  XFER %d '%s'", amount, color.String())
		transfer[color] = amount
		return true
	})
	return transfer
}

func (o *SoloScContext) postSync(contract, function iscp.Hname, paramsID, transferID, delay int32) {
	if delay != 0 {
		o.Panicf("unsupported nonzero delay for SoloContext")
	}
	ctx := o.ctx
	if contract != iscp.Hn(ctx.scName) {
		o.Panicf("invalid contract")
	}
	funcName := ctx.wc.FunctionFromCode(uint32(function))
	if funcName == "" {
		o.Panicf("unknown function")
	}
	o.Tracef("POST %s.%s", ctx.scName, funcName)
	params := o.getParams(paramsID)
	req := solo.NewCallParamsFromDic(ctx.scName, funcName, params)
	if transferID != 0 {
		transfer := o.getTransfer(transferID)
		req.WithTransfers(transfer)
	}
	if ctx.mint > 0 {
		mintAddress := ledgerstate.NewED25519Address(ctx.keyPair.PublicKey)
		req.WithMint(mintAddress, ctx.mint)
	}
	_ = wasmhost.Connect(ctx.wasmHostOld)
	var res dict.Dict
	if ctx.offLedger {
		ctx.offLedger = false
		res, ctx.Err = ctx.Chain.PostRequestOffLedger(req, ctx.keyPair)
	} else if !ctx.isRequest {
		ctx.Tx, res, ctx.Err = ctx.Chain.PostRequestSyncTx(req, ctx.keyPair)
	} else {
		ctx.isRequest = false
		ctx.Tx, _, ctx.Err = ctx.Chain.RequestFromParamsToLedger(req, nil)
		if ctx.Err == nil {
			ctx.Chain.Env.EnqueueRequests(ctx.Tx)
		}
	}
	_ = wasmhost.Connect(ctx.wc)
	if ctx.Err != nil {
		return
	}
	returnID := o.GetObjectID(wasmhost.KeyReturn, wasmhost.OBJTYPE_MAP)
	ctx.wc.FindObject(returnID).(*wasmproc.ScDict).SetKvStore(res)
}
