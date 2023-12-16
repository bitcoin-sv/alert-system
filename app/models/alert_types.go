package models

// AlertType is the type of alert
type AlertType uint32

// AlertTypeInformational an alert type for informational alerts
const AlertTypeInformational AlertType = 0x01

// AlertTypeFreezeUtxo is an alert type for freezing a UTXO
const AlertTypeFreezeUtxo AlertType = 0x02

// AlertTypeUnfreezeUtxo is an alert type for unfreezing a UTXO
const AlertTypeUnfreezeUtxo AlertType = 0x03

// AlertTypeConfiscateUtxo is an alert type for confiscating a UTXO
const AlertTypeConfiscateUtxo AlertType = 0x04

// AlertTypeBanPeer is an alert type for banning a peer
const AlertTypeBanPeer AlertType = 0x05

// AlertTypeUnbanPeer is an alert type for unbanning a peer
const AlertTypeUnbanPeer AlertType = 0x06

// AlertTypeInvalidateBlock is an alert type for invalidating a block
const AlertTypeInvalidateBlock AlertType = 0x07

// AlertTypeSetKeys is an alert type for setting keys
const AlertTypeSetKeys AlertType = 0x08
