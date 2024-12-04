// Package main is a package for hacks and running tests
package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/bitcoin-sv/alert-system/app/p2p"
	"github.com/bitcoin-sv/alert-system/utils"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	models2 "github.com/libsv/go-bn/models"
)

func main() {
	alertTypeFlag := flag.Uint("type", uint(1), "type of alert to publish")
	sequenceNumber := flag.Uint("sequence", uint(1), "sequence number to publish")
	//pubKeys := flag.String("pub-keys", "", "public keys to be used for set keys")
	blockHash := flag.String("block-hash", "", "block hash to invalidate")
	//peer := flag.String("peer", "", "peer to ban/unban")
	keys := flag.String("signing-keys", "", "signing keys")

	flag.Parse()

	ctx := context.Background()

	// Load the configuration and services
	_appConfig, err := config.LoadDependencies(context.Background(), models.BaseModels, false)
	if err != nil {
		log.Fatalf("error loading configuration: %s", err.Error())
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
		a = InfoAlert(*sequenceNumber, "Testing block invalidation on testnet of 00000000000439a2c310b4e457f7e36f51c25931ccda8d512aeb2300587bcd5d", model.WithAllDependencies(_appConfig))
	case models.AlertTypeInvalidateBlock:

		a = invalidateBlockAlert(*sequenceNumber, *blockHash, model.WithAllDependencies(_appConfig))
	case models.AlertTypeBanPeer:
		//a = BanPeerAlert(*sequenceNumber, *peer)
	case models.AlertTypeUnbanPeer:
		//a = UnbanPeerAlert(*sequenceNumber, *peer)
	case models.AlertTypeConfiscateUtxo:
		a = confiscateAlert(*sequenceNumber, model.WithAllDependencies(_appConfig))
	case models.AlertTypeFreezeUtxo:
		a = freezeAlert(*sequenceNumber, model.WithAllDependencies(_appConfig))
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

	// Create the p2p server
	var p2pServer *p2p.Server
	if p2pServer, err = p2p.NewServer(p2p.ServerOptions{
		TopicNames: []string{_appConfig.P2P.TopicName},
		Config:     _appConfig,
	}); err != nil {
		_appConfig.Services.Log.Fatalf("error creating p2p server: %s", err.Error())
	}

	// Start the p2p server
	if err = p2pServer.Start(context.Background()); err != nil {
		_appConfig.Services.Log.Fatalf("error starting p2p server: %s", err.Error())
	}

	// Wait for server to be connected
	for !p2pServer.Connected() {
		time.Sleep(1 * time.Second)
	}
	topics := p2pServer.Topics()

	// Check if alert signatures are valid
	// This is going to check the local database to see the active keys and ensure signatures come from them
	var v bool
	if v, err = a.AreSignaturesValid(ctx); err != nil {
		panic(err)
	}
	if !v {
		_appConfig.Services.Log.Errorf("signature is not valid")
		return
	}
	_appConfig.Services.Log.Infof("raw alert bytes: %x", a.Serialize())
	// publish the alert
	publish(ctx, topics[_appConfig.P2P.TopicName], a.Serialize())
	_appConfig.Services.Log.Infof("successfully published alert to topic %s", _appConfig.P2P.TopicName)
}

// InfoAlert creates an informational alert
func InfoAlert(seq uint, msg string, opts ...model.Options) *models.AlertMessage {
	// Create the new alert
	opts = append(opts, model.New())
	newAlert := models.NewAlertMessage(opts...)
	newAlert.SetAlertType(models.AlertTypeInformational)
	newAlert.SetRawMessage([]byte(msg))
	newAlert.SequenceNumber = uint32(seq)
	newAlert.SetTimestamp(uint64(time.Now().Second()))
	newAlert.SetVersion(0x01)

	newAlert.SerializeData()
	return newAlert
}

func freezeAlert(seq uint, opts ...model.Options) *models.AlertMessage {
	tx, _ := hex.DecodeString("d83dee7aec89a9437345d9676bc727a2592e5b3988f4343931181f86b666eace")
	fund := models.Fund{
		TransactionOutID:           [32]byte(tx),
		Vout:                       uint64(0),
		EnforceAtHeightStart:       uint64(10000),
		EnforceAtHeightEnd:         uint64(10100),
		PolicyExpiresWithConsensus: false,
	}
	opts = append(opts, model.New())
	newAlert := models.NewAlertMessage(opts...)
	newAlert.SetAlertType(models.AlertTypeFreezeUtxo)
	newAlert.SetRawMessage(fund.Serialize())
	newAlert.SequenceNumber = uint32(seq)
	newAlert.SetTimestamp(uint64(time.Now().Second()))
	newAlert.SetVersion(0x01)
	newAlert.SerializeData()
	return newAlert
}

func confiscateAlert(seq uint, opts ...model.Options) *models.AlertMessage {
	tx := models2.ConfiscationTransactionDetails{
		ConfiscationTransaction: models2.ConfiscationTransaction{
			Hex:             "dd1b08331cf22da4d27bd1b29019a04a168805d49b48d65a7fec381eb4307d61",
			EnforceAtHeight: 10000,
		},
	}
	raw := []byte{}
	enforce := [8]byte{}
	binary.LittleEndian.PutUint64(enforce[:], uint64(tx.ConfiscationTransaction.EnforceAtHeight))
	raw = append(raw, enforce[:]...)
	by, _ := hex.DecodeString(tx.ConfiscationTransaction.Hex)
	raw = append(raw, by...)
	opts = append(opts, model.New())
	newAlert := models.NewAlertMessage(opts...)
	newAlert.SetAlertType(models.AlertTypeConfiscateUtxo)
	newAlert.SetRawMessage(raw)
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
*/
// invalidateBlockAlert creates an invalidate block alert
func invalidateBlockAlert(seq uint, blockHash string, opts ...model.Options) *models.AlertMessage {
	hash, err := hex.DecodeString(blockHash)
	if err != nil {
		panic(err)
	}
	raw := []byte{}

	opts = append(opts, model.New())
	raw = append(raw, hash...)
	raw = append(raw, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6E, 0x67)
	newAlert := models.NewAlertMessage(opts...)
	newAlert.SetAlertType(models.AlertTypeInvalidateBlock)
	newAlert.SetVersion(0x01)
	newAlert.SetTimestamp(uint64(time.Now().Unix()))
	newAlert.SequenceNumber = uint32(seq)
	newAlert.SetRawMessage(raw)
	newAlert.SerializeData()

	return newAlert
}

/*
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
	// Sleep for a second just to ensure it was fully processed
	time.Sleep(1 * time.Second)
}
