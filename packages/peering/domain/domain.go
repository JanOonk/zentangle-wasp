package domain

import (
	"crypto/rand"
	"sort"
	"sync"
	"time"

	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
)

type DomainImpl struct {
	attachedTo  peering.PeeringID
	netProvider peering.NetworkProvider
	nodes       map[string]peering.PeerSender
	permutation *util.Permutation16
	netIDs      []string
	log         *logger.Logger
	mutex       *sync.RWMutex
}

// NewPeerDomain creates a collection. Ignores self
func NewPeerDomain(netProvider peering.NetworkProvider, initialNodes []peering.PeerSender, log *logger.Logger) *DomainImpl {
	ret := &DomainImpl{
		netProvider: netProvider,
		nodes:       make(map[string]peering.PeerSender),
		permutation: util.NewPermutation16(uint16(len(initialNodes)), nil),
		netIDs:      make([]string, 0, len(initialNodes)),
		log:         log,
		mutex:       &sync.RWMutex{},
	}
	for _, sender := range initialNodes {
		ret.nodes[sender.NetID()] = sender
	}
	ret.reshufflePeers()
	return ret
}

func NewPeerDomainByNetIDs(netProvider peering.NetworkProvider, peerNetIDs []string, log *logger.Logger) (*DomainImpl, error) {
	peers := make([]peering.PeerSender, 0, len(peerNetIDs))
	for _, nid := range peerNetIDs {
		if nid == netProvider.Self().NetID() {
			continue
		}
		peer, err := netProvider.PeerByNetID(nid)
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}
	return NewPeerDomain(netProvider, peers, log), nil
}

func (d *DomainImpl) SendMsgByNetID(netID string, msg *peering.PeerMessage) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	peer, ok := d.nodes[netID]
	if !ok {
		d.log.Warnf("SendMsgByNetID: NetID %v is not in the domain", netID)
		return
	}
	peer.SendMsg(msg)
}

func (d *DomainImpl) SendSimple(netID string, msgType byte, msgData []byte) {
	d.SendMsgByNetID(netID, &peering.PeerMessage{
		PeeringID: d.attachedTo,
		Timestamp: time.Now().UnixNano(),
		MsgType:   msgType,
		MsgData:   msgData,
	})
}

func (d *DomainImpl) SendMsgToRandomPeers(upToNumPeers uint16, msg *peering.PeerMessage) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if int(upToNumPeers) > len(d.nodes) {
		upToNumPeers = uint16(len(d.nodes))
	}
	for i := uint16(0); i < upToNumPeers; i++ {
		d.SendMsgByNetID(d.netIDs[d.permutation.Next()], msg)
	}
}

func (d *DomainImpl) SendMsgToRandomPeersSimple(upToNumPeers uint16, msgType byte, msgData []byte) {
	d.SendMsgToRandomPeers(upToNumPeers, &peering.PeerMessage{
		PeeringID: d.attachedTo,
		Timestamp: time.Now().UnixNano(),
		MsgType:   msgType,
		MsgData:   msgData,
	})
}

func (d *DomainImpl) GetRandomPeers(upToNumPeers int) []string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if upToNumPeers > len(d.netIDs) {
		upToNumPeers = len(d.netIDs)
	}
	ret := make([]string, upToNumPeers)
	for i := range ret {
		ret[i] = d.netIDs[d.permutation.Next()]
	}
	return ret
}

func (d *DomainImpl) AddPeer(netID string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, ok := d.nodes[netID]; ok {
		return nil
	}
	if netID == d.netProvider.Self().NetID() {
		return nil
	}
	peer, err := d.netProvider.PeerByNetID(netID)
	if err != nil {
		return err
	}
	d.nodes[netID] = peer
	d.reshufflePeers()

	return nil
}

func (d *DomainImpl) RemovePeer(netID string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.nodes, netID)
	d.reshufflePeers()
}

func (d *DomainImpl) ReshufflePeers(seedBytes ...[]byte) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.reshufflePeers(seedBytes...)
}

func (d *DomainImpl) reshufflePeers(seedBytes ...[]byte) {
	d.netIDs = make([]string, 0, len(d.nodes))
	for netID := range d.nodes {
		d.netIDs = append(d.netIDs, netID)
	}
	sort.Strings(d.netIDs)
	var seedB []byte
	if len(seedBytes) == 0 {
		var b [8]byte
		seedB = b[:]
		_, _ = rand.Read(seedB)
	} else {
		seedB = seedBytes[0]
	}
	d.permutation.Shuffle(seedB)
}

func (d *DomainImpl) Attach(peeringID *peering.PeeringID, callback func(recv *peering.RecvEvent)) interface{} {
	d.attachedTo = *peeringID
	return d.netProvider.Attach(peeringID, func(recv *peering.RecvEvent) {
		peer, ok := d.nodes[recv.From.NetID()]
		if ok && peer.NetID() != d.netProvider.Self().NetID() {
			recv.Msg.SenderNetID = peer.NetID()
			callback(recv)
			return
		}
		d.log.Warnf("dropping message MsgType=%v from %v, it does not belong to the peer domain.",
			recv.Msg.MsgType, recv.From.NetID())
	})
}

func (d *DomainImpl) Detach(attachID interface{}) {
	d.netProvider.Detach(attachID)
}

func (d *DomainImpl) Close() {
	for i := range d.nodes {
		d.nodes[i].Close()
	}
}
