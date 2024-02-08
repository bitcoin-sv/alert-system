package p2p

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-libp2p/core/discovery"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/bitcoin-sv/alert-system/app/webhook"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/mrz1836/go-datastore"
)

// Define an interface to handle topic notifications
// TODO Likely need to come up with a more standard way to support this with
// multiple topics. But this allows an external service to use this package and
// handle subscription events

// ServerOptions are the options for the server
type ServerOptions struct {
	Config     *config.Config
	TopicNames []string
}

// Server is the P2P server
type Server struct {
	// alertKeyTopicName string
	connected                     bool
	config                        *config.Config
	host                          host.Host
	privateKey                    *crypto.PrivKey
	subscriptions                 map[string]*pubsub.Subscription
	topicNames                    []string
	topics                        map[string]*pubsub.Topic
	dht                           *dht.IpfsDHT
	quitAlertProcessingChannel    chan bool
	quitPeerDiscoveryChannel      chan bool
	quitPeerInitializationChannel chan bool
	//peers         []peer.AddrInfo
}

// NewServer will create a new server
// Instantiate a new server instance, optionally include a subscriber
// if `subscriber` is nil, we won't process the subscription events
func NewServer(o ServerOptions) (*Server, error) {

	o.Config.Services.Log.Debug("creating P2P service")

	// Attempt to read the private key from the file
	pk, err := readPrivateKey(o.Config.P2P.PrivateKeyPath)
	if err != nil {

		// If the file doesn't exist, generate a new private key
		if pk, err = generatePrivateKey(o.Config.P2P.PrivateKeyPath); err != nil {
			return nil, err
		}
	}

	// Create a new host
	var h host.Host
	if h, err = libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%s", o.Config.P2P.IP, o.Config.P2P.Port)),
		libp2p.Identity(*pk),
		libp2p.EnableHolePunching(),
	); err != nil {
		return nil, err
	}

	// Print out the peer ID and addresses
	o.Config.Services.Log.Debugf("peer ID: %s", h.ID().String())
	o.Config.Services.Log.Info("connect to me on:")
	for _, addr := range h.Addrs() {
		o.Config.Services.Log.Infof(" %s/p2p/%s", addr, h.ID().String())
	}

	// Return the server
	return &Server{
		host:                          h,
		topicNames:                    o.TopicNames,
		privateKey:                    pk,
		config:                        o.Config,
		quitPeerInitializationChannel: make(chan bool),
	}, nil
}

// Start the server and subscribe to all topics
func (s *Server) Start(ctx context.Context) error {
	s.config.Services.Log.Info("p2p service initializing & starting")

	// Initialize the DHT
	kademliaDHT, err := s.initDHT(ctx)
	if err != nil {
		return err
	}
	s.dht = kademliaDHT

	// Advertise our existence so that other peers can find us
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	for _, topicName := range s.topicNames {
		dutil.Advertise(ctx, routingDiscovery, topicName)
	}

	s.quitPeerDiscoveryChannel = s.RunPeerDiscovery(ctx, routingDiscovery)
	s.quitAlertProcessingChannel = s.RunAlertProcessingCron(ctx)

	ps, err := pubsub.NewGossipSub(ctx, s.host, pubsub.WithDiscovery(routingDiscovery))
	if err != nil {
		return err
	}
	topics := map[string]*pubsub.Topic{}
	subscriptions := map[string]*pubsub.Subscription{}

	s.host.SetStreamHandler(protocol.ID(s.config.P2P.AlertSystemProtocolID), func(stream network.Stream) {
		t := StreamThread{
			stream: stream,
			config: s.config,
			ctx:    ctx,
			peer:   stream.Conn().RemotePeer(),
		}

		if err = t.ProcessSyncMessage(ctx); err != nil {
			s.config.Services.Log.Errorf("failed to process sync message: %v", err.Error())
			//_ = stream.Reset()
		} else {
			s.config.Services.Log.Debugf("closing stream %v for peer %v", stream.ID(), t.peer.String())
			//_ = stream.Close()
		}
		//_ = stream.Close()
	})

	s.config.Services.Log.Debugf("stream handler set")
	for !s.connected {
		time.Sleep(5 * time.Second)
	}
	for _, topicName := range s.topicNames {
		var topic *pubsub.Topic
		if topic, err = ps.Join(topicName); err != nil {
			return err
		}
		topics[topicName] = topic

		var sub *pubsub.Subscription
		if sub, err = topic.Subscribe(); err != nil {
			return err
		}
		subscriptions[topicName] = sub

		// Sync the subscriber
		go s.Subscribe(ctx, sub, s.host.ID())
	}
	s.topics = topics
	s.subscriptions = subscriptions
	s.config.Services.Log.Infof("P2P successfully started")
	go func() {
		for { //nolint:gosimple // This is the only way to perform this loop at the moment
			select {
			case <-ctx.Done():
				s.config.Services.Log.Info("p2p service shutting down")
				return
			}
		}
	}()
	return nil
}

// Connected returns true if the server is connected
func (s *Server) Connected() bool {
	return s.connected
}

// Stop the server
func (s *Server) Stop(_ context.Context) error {
	// todo there needs to be a way to stop the server
	s.config.Services.Log.Info("stopping P2P service")
	s.quitPeerDiscoveryChannel <- true
	s.quitAlertProcessingChannel <- true
	s.quitPeerInitializationChannel <- true
	return nil
}

// RunAlertProcessingCron starts a cron job to attempt to retry unprocessed alerts
func (s *Server) RunAlertProcessingCron(ctx context.Context) chan bool {
	ticker := time.NewTicker(s.config.AlertProcessingInterval)
	quit := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.processAlerts(ctx)
				if err != nil {
					s.config.Services.Log.Errorf("error processing alerts: %v", err.Error())
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}

// processAlerts performs the alert processing
func (s *Server) processAlerts(ctx context.Context) error {
	alerts, err := models.GetAllUnprocessedAlerts(ctx, nil, model.WithAllDependencies(s.config))
	if err != nil {
		return err
	}
	s.config.Services.Log.Infof("Attempting to process %d failed alerts", len(alerts))
	success := 0
	for _, alert := range alerts {
		alert.SetOptions(model.WithAllDependencies(s.config))
		// Serialize the alert data and hash
		err := alert.ReadRaw()
		if err != nil {
			continue
		}
		alert.SerializeData()
		// Process the alert
		ak := alert.ProcessAlertMessage()
		if ak == nil {
			continue
		}
		if err = ak.Read(alert.GetRawMessage()); err != nil {
			return err
		}
		s.config.Services.Log.Debugf("attempting to process alert %d of type %d", alert.SequenceNumber, alert.GetAlertType())
		alert.Processed = true
		if err = ak.Do(ctx); err != nil {
			s.config.Services.Log.Errorf("failed to process alert %d; err: %v", alert.SequenceNumber, err.Error())
			alert.Processed = false
		}

		if alert.Processed {
			success++
			// Save the alert
			if err = alert.Save(ctx); err != nil {
				return err
			}
		}
	}
	s.config.Services.Log.Infof("Processed %d failed alerts", success)
	return nil
}

// RunPeerDiscovery starts a cron job to resync peers and update routable peers
func (s *Server) RunPeerDiscovery(ctx context.Context, routingDiscovery *drouting.RoutingDiscovery) chan bool {
	ticker := time.NewTicker(s.config.P2P.PeerDiscoveryInterval)
	quit := make(chan bool, 1)
	go func() {
		err := s.discoverPeers(ctx, routingDiscovery)
		if err != nil {
			s.config.Services.Log.Errorf("error discovering peers: %v", err.Error())
		}
		for {
			select {
			case <-ticker.C:
				err := s.discoverPeers(ctx, routingDiscovery)
				if err != nil {
					s.config.Services.Log.Errorf("error discovering peers: %v", err.Error())
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}

// generatePrivateKey generates a private key and stores it in `private_key` file
func generatePrivateKey(filePath string) (*crypto.PrivKey, error) {
	// Generate a new key pair
	privateKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Convert private key to bytes
	var privateBytes []byte
	if privateBytes, err = crypto.MarshalPrivateKey(privateKey); err != nil {
		return nil, err
	}

	// Save private key to a file
	if err = os.WriteFile(filePath, privateBytes, 0644); err != nil { //nolint:gosec // This is a local private key
		return nil, err
	}

	return &privateKey, nil
}

// readPrivateKey reads a private key from `private_key` file
func readPrivateKey(filePath string) (*crypto.PrivKey, error) {
	// Read private key from a file
	privateBytes, err := os.ReadFile(filePath) //nolint:gosec // This is a local private key
	if err != nil {
		return nil, err
	}

	// Unmarshal the private key bytes into a key
	var privateKey crypto.PrivKey
	if privateKey, err = crypto.UnmarshalPrivateKey(privateBytes); err != nil {
		return nil, err
	}

	return &privateKey, nil
}

// Subscriptions lists all current subscriptions
func (s *Server) Subscriptions() map[string]*pubsub.Subscription {
	return s.subscriptions
}

// Topics lists all topics
func (s *Server) Topics() map[string]*pubsub.Topic {
	return s.topics
}

// discoverPeers will discover peers
func (s *Server) discoverPeers(ctx context.Context, routingDiscovery *drouting.RoutingDiscovery) error {
	s.config.Services.Log.Infof("Running peer discovery at %s", time.Now().String())

	// Look for others who have announced and attempt to connect to them
	connected := 0
	for connected < 2 {
		for _, topicName := range s.topicNames {
			s.config.Services.Log.Debugf("searching for peers for topic %s..\n", topicName)

			var peerChan <-chan peer.AddrInfo
			var err error
			if peerChan, err = routingDiscovery.FindPeers(ctx, topicName, discovery.TTL(1*time.Minute)); err != nil {
				return err
			}

			// Loop through all peers found
			for foundPeer := range peerChan {

				// Don't connect to ourselves
				if foundPeer.ID == s.host.ID() {
					continue // No self connection
				}

				// Failed to connect to peer
				s.config.Services.Log.Debugf("attempting connection to %s", foundPeer.ID.String())

				if err = s.host.Connect(ctx, foundPeer); err != nil {
					// we fail to connect to a lot of peers. Just ignore it for now.
					s.config.Services.Log.Debugf("failed connecting to %s, error: %s", foundPeer.ID.String(), err.Error())
					continue
				}

				// Connected to peer
				s.config.Services.Log.Infof("connected to: %s", foundPeer.ID.String())

				// Open a stream to the peer
				var stream network.Stream
				if stream, err = s.host.NewStream(ctx, foundPeer.ID, protocol.ID(s.config.P2P.AlertSystemProtocolID)); err != nil {
					s.config.Services.Log.Debugf("failed new stream to %s error: %s", foundPeer.ID.String(), err.Error())
					continue
				}

				// Sync the stream thread
				t := StreamThread{
					config:      s.config,
					ctx:         ctx,
					peer:        foundPeer.ID,
					stream:      stream,
					quitChannel: s.quitPeerDiscoveryChannel,
				}

				// Sync the stream thread
				if err = t.Sync(ctx); err != nil {
					s.config.Services.Log.Debugf("failed to start stream thread to %s error: %s", foundPeer.ID.String(), err.Error())
					continue
				}

				s.config.Services.Log.Infof("successfully synced up to %d from peer %s", t.LatestSequence(), foundPeer.ID.String())

				// Set the flag
				connected++
			}
			time.Sleep(1 * time.Second)
		}
	}

	// We are connected
	s.config.Services.Log.Debugf("peer discovery complete")
	s.config.Services.Log.Debugf("connected to %d peers\n", len(s.host.Network().Peers()))
	s.config.Services.Log.Debugf("peerstore has %d peers\n", len(s.host.Peerstore().Peers()))
	s.config.Services.Log.Infof("Successfully discovered %d active peers at %s", connected, time.Now().String())
	s.connected = true
	return nil
}

// Subscribe will subscribe to the alert system
func (s *Server) Subscribe(ctx context.Context, subscriber *pubsub.Subscription, hostID peer.ID) {
	s.config.Services.Log.Infof("subscribing to %s topic", subscriber.Topic())
	for {

		msg, err := subscriber.Next(ctx)

		if err != nil {
			s.config.Services.Log.Infof("error subscribing via next: %s", err.Error())
			continue
		}

		// only consider messages delivered by other peers
		if msg.ReceivedFrom == hostID {
			continue
		}

		// Read the alert key header
		var ak *models.AlertMessage
		if ak, err = models.NewAlertFromBytes(msg.Data, model.WithAllDependencies(s.config)); err != nil {
			s.config.Services.Log.Errorf("error reading alert key: %s", err.Error())
			continue
		}

		// Set the hash
		ak.SerializeData()

		// Ensure signatures are valid
		var valid bool
		if valid, err = ak.AreSignaturesValid(ctx); err != nil {
			s.config.Services.Log.Infof("error verifying signatures: %s", err.Error())
			continue
		}

		// Ensure the signature is valid
		if !valid {
			// TODO save these messages still and ban the peer?
			s.config.Services.Log.Info("signature block is invalid")
			continue
		}

		// Ensure the sequence number is correct
		if _, err = models.GetAlertMessageBySequenceNumber(
			ctx, ak.SequenceNumber-1, model.WithAllDependencies(s.config),
		); err != nil {
			// TODO save these messages still and ban the peer? and possibly resync
			s.config.Services.Log.Errorf("failed to find prior sequenced alert (num %d): %s", ak.SequenceNumber-1, err.Error())
			continue
		}

		// Check if the alert already exists
		var dup *models.AlertMessage
		if dup, err = models.GetAlertMessageBySequenceNumber(
			ctx, ak.SequenceNumber, model.WithAllDependencies(s.config),
		); err == nil && dup != nil && len(dup.Hash) > 0 {
			// TODO save these messages still?
			s.config.Services.Log.Errorf("alert %s already has sequence number %d", dup.Hash, ak.SequenceNumber)
			continue
		}

		// Did we get a real error?
		if err != nil && !errors.Is(err, datastore.ErrNoResults) {
			s.config.Services.Log.Errorf("error looking for duplicate alert: %s", err.Error())
			continue
		}

		// Process the alert message into correct interface
		am := ak.ProcessAlertMessage()
		if err = am.Read(ak.GetRawMessage()); err != nil {
			s.config.Services.Log.Errorf("failed to read message: %s", err.Error())
			continue
		}
		ak.Processed = true

		// Perform alert action
		if err = am.Do(ctx); err != nil {
			s.config.Services.Log.Errorf("failed to do alert action: %s", err.Error())
			ak.Processed = false
		}

		// Save the alert message
		if err = ak.Save(ctx); err != nil {
			s.config.Services.Log.Errorf("failed to save alert message: %s", err.Error())
		}

		s.config.Services.Log.Infof("[%s] got alert type: %d, from: %s", subscriber.Topic(), ak.GetAlertType(), msg.ReceivedFrom.String())

		// Send the webhook
		if len(s.config.AlertWebhookURL) > 0 {
			if err = webhook.PostAlert(ctx, s.config.Services.HTTPClient, s.config.AlertWebhookURL, ak); err != nil {
				s.config.Services.Log.Errorf("error processing webhook request: %s", err.Error())
			}
		}
	}
}
