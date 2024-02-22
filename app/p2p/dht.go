package p2p

import (
	"context"
	"sync"
	"time"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/multiformats/go-multiaddr"

	dht "github.com/libp2p/go-libp2p-kad-dht"
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
	options = append(options, dht.QueryFilter(dht.PublicQueryFilter))

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
	if s.config.P2P.BootstrapPeer != "" {
		// Connect to the chosen ipfs nodes
		var pubPeer multiaddr.Multiaddr
		if pubPeer, err = multiaddr.NewMultiaddr(s.config.P2P.BootstrapPeer); err != nil {
			return nil, err
		}
		peers = append(peers, pubPeer)
	}

	// Connect to the chosen ipfs nodes
	connected := false
	for !connected {
		select {
		case <-s.quitPeerInitializationChannel:
			return kademliaDHT, nil
		default:
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
	}

	return kademliaDHT, nil
}
