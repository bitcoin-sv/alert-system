// Package main is a package for hacks and running tests
package main

/*func main() {
	alertTypeFlag := flag.Uint("type", uint(1), "type of alert to publish")
	sequenceNumber := flag.Uint("sequence", uint(1), "sequence number to publish")
	pubKeys := flag.String("pub-keys", "", "public keys to be used for set keys")
	blockHash := flag.String("block-hash", "", "block hash to invalidate")
	peer := flag.String("peer", "", "peer to ban/unban")
	keys := flag.String("signing-keys", "", "signing keys")

	flag.Parse()

	ctx := context.Background()
	log := gocore.Log("sig-test")
	c, err := database.NewSqliteClient(ctx, log)
	if err != nil {
		panic(err)
	}

	if err = util.EnsureGenesis(ctx, c); err != nil {
		panic(err)
	}

	alertType := alert.MessageType(*alertTypeFlag)
	a := alert.Alert{}
	switch alertType {
	case alert.MessageTypeInformational:
		a = InfoAlert(*sequenceNumber)
	case alert.MessageTypeInvalidateBlock:
		a = InvalidateBlockAlert(*sequenceNumber, *blockHash)
	case alert.MessageTypeBanPeer:
		a = BanPeerAlert(*sequenceNumber, *peer)
	case alert.MessageTypeUnbanPeer:
		a = UnbanPeerAlert(*sequenceNumber, *peer)
	case alert.MessageTypeConfiscateUtxo:
		panic(fmt.Errorf("not implemented"))
	case alert.MessageTypeFreezeUtxo:
		panic(fmt.Errorf("not implemented"))
	case alert.MessageTypeUnfreezeUtxo:
		panic(fmt.Errorf("not implemented"))
	case alert.MessageTypeSetKeys:
		publicKeys := strings.Split(*pubKeys, ",")
		if len(publicKeys) != 5 {
			panic(fmt.Errorf("did not get 5 public keys to set"))
		}
		a = SetKeys(*sequenceNumber, publicKeys)
	}
	a.Datastore = c
	a.Log = log

	var sigs [][]byte
	if *keys == "" {
		if sigs, err = utils.SignWithGenesis(a.Data); err != nil {
			panic(err)
		}
	} else {
		privKeys := strings.Split(*keys, ",")
		if len(privKeys) != 3 {
			panic(fmt.Errorf("3 private keys not supplied"))
		}
		if sigs, err = utils.SignWithKeys(a.Data, privKeys); err != nil {
			panic(err)
		}
	}

	a.Signatures = sigs

	var raw []byte
	if raw, err = a.Serialize(); err != nil {
		panic(err)
	}
	for _, sig := range a.Signatures {
		log.Infof("sig: %x", sig)
	}

	var v bool
	if v, err = a.AreSignaturesValid(ctx); err != nil {
		panic(err)
	}
	if !v {
		log.Errorf("signature is not valid")
		return
	}
	for _, sig := range a.Signatures {
		log.Infof("signature: %x", sig)
	}
	log.Infof("data: %x", a.Data)
	log.Infof("alert: %#v", a)
	log.Infof("raw: %x", raw)

}

// InfoAlert creates an informational alert
func InfoAlert(seq uint) alert.Alert {
	a := alert.Alert{
		Version:        0x01,
		SequenceNumber: uint32(seq),
		Timestamp:      uint64(time.Now().Second()),
		AlertType:      0x01,
		AlertMessage:   []byte{0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
	}
	data, err := a.SerializeData()
	if err != nil {
		panic(err)
	}
	a.Data = data
	return a
}

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
	var msg []byte
	msg = append(msg, hash...)
	msg = append(msg, []byte{0x01, 0x01}...) // Just append a 1 byte reason for simplicity
	a := alert.Alert{
		Version:        0x01,
		SequenceNumber: uint32(seq),
		Timestamp:      uint64(time.Now().Second()),
		AlertType:      alert.MessageTypeInvalidateBlock,
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
		AlertType:      alert.MessageTypeSetKeys,
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
}
*/
