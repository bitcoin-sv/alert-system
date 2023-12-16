// Package main is a package for hacks and running tests
package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bitcoin-sv/alert-system/app/models/model"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/bitcoin-sv/alert-system/app/p2p"
	"github.com/bitcoin-sv/alert-system/utils"
	"github.com/ordishs/gocore"
)

func main() {
	alertTypeFlag := flag.Uint("type", uint(1), "type of alert to publish")
	sequenceNumber := flag.Uint("sequence", uint(1), "sequence number to publish")
	//pubKeys := flag.String("pub-keys", "", "public keys to be used for set keys")
	//blockHash := flag.String("block-hash", "", "block hash to invalidate")
	//peer := flag.String("peer", "", "peer to ban/unban")
	keys := flag.String("signing-keys", "", "signing keys")

	flag.Parse()

	ctx := context.Background()
	log := gocore.Log("sig-test")

	// Load the configuration and services
	_appConfig, err := config.LoadConfig(context.Background(), models.BaseModels, false)
	if err != nil {
		_appConfig.Services.Log.Fatalf("error loading configuration: %s", err.Error())
	}
	defer func() {
		_appConfig.CloseAll(context.Background())
	}()
	err = models.CreateGenesisAlert(
		context.Background(), model.WithAllDependencies(_appConfig),
	)
	if err != nil {
		panic(err)
	}

	alertType := models.AlertType(*alertTypeFlag)
	a := &models.AlertMessage{}
	switch alertType {
	case models.AlertTypeInformational:
		a = InfoAlert(*sequenceNumber, model.WithAllDependencies(_appConfig))
	case models.AlertTypeInvalidateBlock:
		//a = InvalidateBlockAlert(*sequenceNumber, *blockHash)
	case models.AlertTypeBanPeer:
		//a = BanPeerAlert(*sequenceNumber, *peer)
	case models.AlertTypeUnbanPeer:
		//a = UnbanPeerAlert(*sequenceNumber, *peer)
	case models.AlertTypeConfiscateUtxo:
		panic(fmt.Errorf("not implemented"))
	case models.AlertTypeFreezeUtxo:
		panic(fmt.Errorf("not implemented"))
	case models.AlertTypeUnfreezeUtxo:
		panic(fmt.Errorf("not implemented"))
	case models.AlertTypeSetKeys:
		//publicKeys := strings.Split(*pubKeys, ",")
		//if len(publicKeys) != 5 {
		//	panic(fmt.Errorf("did not get 5 public keys to set"))
		//}
		//a = SetKeys(*sequenceNumber, publicKeys)
	}

	var sigs [][]byte
	if *keys == "" {
		if sigs, err = utils.SignWithGenesis(a.GetRawData()); err != nil {
			panic(err)
		}
	} else {
		privKeys := strings.Split(*keys, ",")
		if len(privKeys) != 3 {
			panic(fmt.Errorf("3 private keys not supplied"))
		}
		if sigs, err = utils.SignWithKeys(a.GetRawData(), privKeys); err != nil {
			panic(err)
		}
	}

	a.SetSignatures(sigs)

	var v bool
	if v, err = a.AreSignaturesValid(ctx); err != nil {
		panic(err)
	}
	if !v {
		log.Errorf("signature is not valid")
		return
	}

	// Create the p2p server
	var p2pServer *p2p.Server
	if p2pServer, err = p2p.NewServer(p2p.ServerOptions{
		TopicNames: []string{config.DatabasePrefix},
		Config:     _appConfig,
	}); err != nil {
		_appConfig.Services.Log.Fatalf("error creating p2p server: %s", err.Error())
	}

	// Start the p2p server
	if err = p2pServer.Start(context.Background()); err != nil {
		_appConfig.Services.Log.Fatalf("error starting p2p server: %s", err.Error())
	}

	for !p2pServer.Connected() {
		time.Sleep(1 * time.Second)
	}
	topics := p2pServer.Topics()

	log.Infof("bytes: %x", a.Serialize())
	publish(ctx, topics[config.DatabasePrefix], a.Serialize())
	log.Infof("successfully published alert")
}

// InfoAlert creates an informational alert
func InfoAlert(seq uint, opts ...model.Options) *models.AlertMessage {
	// Create the new alert
	opts = append(opts, model.New())
	newAlert := models.NewAlertMessage(opts...)
	newAlert.SetAlertType(models.AlertTypeInformational)
	newAlert.SetRawMessage([]byte{0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f})
	newAlert.SequenceNumber = uint32(seq)
	newAlert.SetTimestamp(uint64(time.Now().Second()))
	newAlert.SetVersion(0x01)

	newAlert.SerializeData()
	return newAlert
}

/*
// BanPeerAlert creates a ban peer alert
func BanPeerAlert(seq uint, peer string) alert.Alert {
	var msg []byte
	peerBytes := []byte(peer)
	buf := bytes.NewBuffer(msg)
	err := wire.WriteVarBytes(buf, 0, peerBytes)
	if err != nil {
		panic(err)
	}
	// Just write 1 byte for reason hack
	buf.WriteByte(0x01)
	buf.WriteByte(0x01)

	a := alert.Alert{
		Version:        0x01,
		SequenceNumber: uint32(seq),
		Timestamp:      uint64(time.Now().Second()),
		AlertType:      0x05,
		AlertMessage:   buf.Bytes(),
	}
	var data []byte
	if data, err = a.SerializeData(); err != nil {
		panic(err)
	}
	a.Data = data
	return a
}

// UnbanPeerAlert creates an unban peer alert
func UnbanPeerAlert(seq uint, peer string) alert.Alert {
	var msg []byte
	peerBytes := []byte(peer)
	buf := bytes.NewBuffer(msg)
	err := wire.WriteVarBytes(buf, 0, peerBytes)
	if err != nil {
		panic(err)
	}
	// Just write 1 byte for reason hack
	buf.WriteByte(0x01)
	buf.WriteByte(0x01)

	a := alert.Alert{
		Version:        0x01,
		SequenceNumber: uint32(seq),
		Timestamp:      uint64(time.Now().Second()),
		AlertType:      0x06,
		AlertMessage:   buf.Bytes(),
	}
	var data []byte
	if data, err = a.SerializeData(); err != nil {
		panic(err)
	}
	a.Data = data
	return a
}

// InvalidateBlockAlert creates an invalidate block alert
func InvalidateBlockAlert(seq uint, blockHash string) alert.Alert {
	hash, err := hex.DecodeString(blockHash)
	if err != nil {
		panic(err)
	}
	msg := []byte{}
	msg = append(msg, hash...)
	msg = append(msg, []byte{0x01, 0x01}...) // Just append a 1 byte reason for simplicity
	a := alert.Alert{
		Version:        0x01,
		SequenceNumber: uint32(seq),
		Timestamp:      uint64(time.Now().Second()),
		AlertType:      models.AlertTypeInvalidateBlock,
		AlertMessage:   msg,
	}
	var data []byte
	if data, err = a.SerializeData(); err != nil {
		panic(err)
	}
	a.Data = data
	return a
}

// SetKeys creates a set keys alert
func SetKeys(seq uint, keys []string) alert.Alert {
	a := alert.Alert{
		Version:        0x01,
		SequenceNumber: uint32(seq),
		Timestamp:      uint64(time.Now().Second()),
		AlertType:      alert.AlertTypeSetKeys,
	}
	var msg []byte
	for _, key := range keys {
		b, _ := hex.DecodeString(key)
		msg = append(msg, b...)
	}
	a.AlertMessage = msg
	data, err := a.SerializeData()
	if err != nil {
		panic(err)
	}
	a.Data = data
	return a
}*/

// publish will publish the data to the topic
func publish(ctx context.Context, topic *pubsub.Topic, data []byte) {
	if err := topic.Publish(ctx, data); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
}
