package p2p

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bitcoin-sv/alert-system/app/config"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// initDHT will initialize the DHT
// Start a DHT, for use in peer discovery. We can't just make a new DHT
// client because we want each peer to maintain its own local copy of the
// DHT, so that the bootstrapping node of the DHT can go down without
// inhibiting future peer discovery.
func (s *Server) initDHT(ctx context.Context) (*dht.IpfsDHT, error) {
	logger := s.config.Services.Log
	var options []dht.Option
	mode := dht.ModeAutoServer
	if s.config.P2P.DHTMode == "client" {
		mode = dht.ModeClient
	}
	options = append(options, dht.Mode(mode))
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
	connected := uint32(0)
	for atomic.LoadUint32(&connected) == 0 {
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

				// Create a local copy of peerInfo for the goroutine
				pi := *peerInfo
				go func(logger config.LoggerInterface, peerInfo peer.AddrInfo) {
					defer wg.Done()
					if localErr := s.host.Connect(ctx, peerInfo); localErr != nil {
						logger.Errorf("bootstrap warning: %s", localErr.Error())
						return
					}
					logger.Infof("connected to peer %v", peerInfo.ID)
					atomic.StoreUint32(&connected, 1)
				}(logger, pi)
			}
			time.Sleep(1 * time.Second)
			wg.Wait()
		}
	}

	return kademliaDHT, nil
}
