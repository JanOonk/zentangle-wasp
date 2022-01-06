// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmclient

import (
	"fmt"
	"strings"
	"time"

	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/wasp/client"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/colored"
	"github.com/iotaledger/wasp/packages/iscp/request"
	"github.com/iotaledger/wasp/packages/iscp/requestargs"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/subscribe"
	"github.com/mr-tron/base58"
)

type Service struct {
	chainID    *iscp.ChainID
	keyPair    *ed25519.KeyPair
	scHname    iscp.Hname
	waspClient *client.WaspClient
}

func (s *Service) Init(svcClient *ServiceClient, chainID string, scHname uint32, eventHandlers map[string]func([]string)) (err error) {
	s.waspClient = svcClient.waspClient
	s.scHname = iscp.Hname(scHname)
	s.chainID, err = iscp.ChainIDFromString(chainID)
	if err != nil {
		return err
	}
	return s.startEventHandlers(svcClient.eventPort, eventHandlers)
}

func (s *Service) AsClientFunc() ClientFunc {
	return ClientFunc{svc: s}
}

func (s *Service) AsClientView() ClientView {
	return ClientView{svc: s}
}

func (s *Service) CallView(viewName string, args *Arguments) (dict.Dict, error) {
	if args == nil {
		args = &Arguments{}
	}
	return s.waspClient.CallView(s.chainID, s.scHname, viewName, args.args)
}

func (s *Service) PostRequest(hFuncName uint32, args *Arguments, transfer *Transfer, keyPair *ed25519.KeyPair) Request {
	bal, err := makeBalances(transfer)
	if err != nil {
		return Request{err: err}
	}
	reqArgs := requestargs.New()
	if args != nil {
		reqArgs.AddEncodeSimpleMany(args.args)
	}
	req := request.NewOffLedger(s.chainID, s.scHname, iscp.Hname(hFuncName), reqArgs)
	req.WithTransfer(bal)
	req.Sign(keyPair)
	err = s.waspClient.PostOffLedgerRequest(s.chainID, req)
	if err != nil {
		return Request{err: err}
	}
	id := req.ID()
	return Request{id: &id}
}

func (s *Service) SignRequests(keyPair *ed25519.KeyPair) {
	s.keyPair = keyPair
}

func (s *Service) WaitRequest(req Request) error {
	return s.waspClient.WaitUntilRequestProcessed(s.chainID, *req.id, 1*time.Minute)
}

func (s *Service) startEventHandlers(eventPort string, handlers map[string]func([]string)) error {
	chMsg := make(chan []string, 20)
	chDone := make(chan bool)
	err := subscribe.Subscribe(eventPort, chMsg, chDone, true, "")
	if err != nil {
		return err
	}
	go func() {
		for {
			for msgSplit := range chMsg {
				event := strings.Join(msgSplit, " ")
				fmt.Printf("%s\n", event)
				if msgSplit[0] == "vmmsg" {
					msg := strings.Split(msgSplit[3], "|")
					handler, ok := handlers[msg[0]]
					if ok {
						handler(msg[1:])
					}
				}
			}
		}
	}()
	return nil
}

/////////////////////////////////////////////////////////////////

func Base58Decode(s string) []byte {
	res, err := base58.Decode(s)
	if err != nil {
		panic("invalid base58 encoding")
	}
	return res
}

func Base58Encode(b []byte) string {
	return base58.Encode(b)
}

func makeBalances(transfer *Transfer) (colored.Balances, error) {
	cb := colored.NewBalances()
	if transfer != nil {
		for color, amount := range transfer.xfer {
			c, err := colored.ColorFromBase58EncodedString(color)
			if err != nil {
				return nil, err
			}
			cb.Set(c, amount)
		}
	}
	return cb, nil
}
