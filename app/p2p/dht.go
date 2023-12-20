package p2p

import (
	"context"
	"sync"
	"time"

	"github.com/multiformats/go-multiaddr"

	"github.com/bitcoin-sv/alert-system/app/config"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
)

// initDHT will initialize the DHT
// Start a DHT, for use in peer discovery. We can't just make a new DHT
// client because we want each peer to maintain its own local copy of the
// DHT, so that the bootstrapping node of the DHT can go down without
// inhibiting future peer discovery.
func (s *Server) initDHT(ctx context.Context) (*dht.IpfsDHT, error) {
	logger := s.config.Services.Log
	var options []dht.Option
	options = append(options, dht.Mode(dht.ModeAutoServer))

	// Sync a DHT, for use in peer discovery. We can't just make a new DHT
	kademliaDHT, err := dht.New(ctx, s.host, options...)
	if err != nil {
		return nil, err
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	// Append the bootstrap nodes
	peers := dht.DefaultBootstrapPeers
	if s.config.P2PBootstrapPeer != "" {
		// Connect to the chosen ipfs nodes
		var pubPeer multiaddr.Multiaddr
		if pubPeer, err = multiaddr.NewMultiaddr(s.config.P2PBootstrapPeer); err != nil {
			return nil, err
		}
		peers = append(peers, pubPeer)
	}

	// Connect to the chosen ipfs nodes
	var connected = false
	for !connected {
		var wg sync.WaitGroup
		for _, peerAddr := range peers {
			var peerInfo *peer.AddrInfo
			if peerInfo, err = peer.AddrInfoFromP2pAddr(peerAddr); err != nil {
				return nil, err
			}
			wg.Add(1)
			go func(logger config.LoggerInterface) {
				defer wg.Done()
				if err = s.host.Connect(ctx, *peerInfo); err != nil {
					logger.Errorf("bootstrap warning: %s", err.Error())
					return
				}
				logger.Infof("connected to peer %v", peerInfo.ID)
				connected = true
			}(logger)
		}
		time.Sleep(1 * time.Second)
		wg.Wait()
	}

	return kademliaDHT, nil
}
