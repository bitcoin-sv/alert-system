package models

/*// Note this test uses the genesis key set in the utils pkg
func TestAlert_AreSignaturesValid(t *testing.T) {
	type args struct {
		Data       string   // convert hex to bytes
		Signatures []string // convert hex to bytes
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid signature",
			args: args{
				Data: "01000000010000003600000000000000010000000568656c6c6f",
				Signatures: []string{
					"2047ce6f2168bb9234f169ac0595a9c4ddca8082195d973c663d71b76172ba40a64fec1cf85adf6c3cd770a520b4cfced4b8491a9bda00b7e90ddf91dd1102f4bd",
					"1fc86fb653b03eb9e4bd4da5c6cd0ccbb187b3ae460f8283bbfa997c83c9ba4a694bfb8935ddafb3eb92f2ed530ae24aeb87431a43bcbecdaa3a5d1c5b15395820",
					"20116b906e78e0ce66dad1bf7c630b7396a6b37e0e14f6bbd664765276331f4944750068efdb781ea9bbb454105a49177bd42f1907243e0cd0a6bb273f0ea46bc4",
				},
			},
			want: true,
		},
		{
			name: "bad signature should fail",
			args: args{
				Data: "01000000123456789012345678901234567890121234567890123456789012345678901202000000050000000c093132372e302e302e310101",
				Signatures: []string{
					"4c6e4b47434b64712b6c387939464759756e676737384b6941627537572b74345768776c3756394554343649314b674245396d7177544c74426f752f41795a564e635358324659425a3032735a476b78596a78464e733d",
					"20116b906e78e0ce66dad1bf7c630b7396a6b37e0e14f6bbd664765276331f4944750068efdb781ea9bbb454105a49177bd42f1907243e0cd0a6bb273f0ea46bc4",
					"1fc86fb653b03eb9e4bd4da5c6cd0ccbb187b3ae460f8283bbfa997c83c9ba4a694bfb8935ddafb3eb92f2ed530ae24aeb87431a43bcbecdaa3a5d1c5b15395820",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			data, err := hex.DecodeString(tt.args.Data)
			if err != nil {
				t.Errorf("failed to decode data hex string for test: %s", err.Error())
				return
			}
			var sigs [][]byte
			for _, s := range tt.args.Signatures {
				var sig []byte
				if sig, err = hex.DecodeString(s); err != nil {
					t.Errorf("failed to decode sig hex string for test: %s", err.Error())
					return
				}
				sigs = append(sigs, sig)
			}

			a := &Alert{
				Config:     nil,
				Data:       data,
				Signatures: sigs,
			}
			if got, _ := a.AreSignaturesValid(ctx); got != tt.want {
				t.Errorf("AreSignaturesValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlert_ProcessAlertMessage(t *testing.T) {
	type fields struct {
		Data         []byte
		Version      uint32
		HashPrevMsg  *chainhash.Hash
		Timestamp    uint64
		AlertType    MessageType
		AlertMessage []byte
		Signatures   [][]byte
	}
	var tests []struct {
		name   string
		fields fields
		want   Alert
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Alert{
				AlertMessage: tt.fields.AlertMessage,
				AlertType:    tt.fields.AlertType,
				Config:       nil,
				Data:         tt.fields.Data,
				Signatures:   tt.fields.Signatures,
				Timestamp:    tt.fields.Timestamp,
				Version:      tt.fields.Version,
			}
			if got := a.ProcessAlertMessage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessAlertMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAlert(t *testing.T) {
	type args struct {
		ak string
	}
	type want struct {
		signatures []string
		data       string
		alertKey   *Alert
	}
	tests := []struct {
		name    string
		args    args
		wanted  want
		wantErr bool
	}{
		{
			name: "valid set ban alert (no verification of peer data)",
			args: args{
				ak: "0100000001000000250000000000000005000000000101201dfc958f68d2ae3d6c66a5152ae104dabc2d81bd455c15f2f792b90f5dac4d8a74c9a5a495392bb04b71e8f3fc45180cdeb0d1f21216420390fffa0e7aa58de1207e5c0433da87b23db2ecc813a7f04a08ef061fc2fc67325d4a6872289c9138491a4724a688587ec5bdd98ea1f25f94d6cb697cd5386cd09f2136398b9a5290ce2014a74f7a2d8acd6887fe0306863529bc702cbec1ecf21f9f848bf57f07a5f78e57eba311bf2240dcc7734e049dd84516472a8805927d7172fe61f3342b580175",
			},
			wanted: want{
				alertKey: &Alert{
					Version:   0x01,
					Timestamp: 0x25,
					AlertType: MessageTypeBanPeer,
				},
				signatures: []string{
					"201dfc958f68d2ae3d6c66a5152ae104dabc2d81bd455c15f2f792b90f5dac4d8a74c9a5a495392bb04b71e8f3fc45180cdeb0d1f21216420390fffa0e7aa58de1",
					"207e5c0433da87b23db2ecc813a7f04a08ef061fc2fc67325d4a6872289c9138491a4724a688587ec5bdd98ea1f25f94d6cb697cd5386cd09f2136398b9a5290ce",
					"2014a74f7a2d8acd6887fe0306863529bc702cbec1ecf21f9f848bf57f07a5f78e57eba311bf2240dcc7734e049dd84516472a8805927d7172fe61f3342b580175",
				},
				data: "0100000001000000250000000000000005000000000101",
			},
		},
		{
			name: "invalid set ban alert",
			args: args{
				ak: "0100000001000000250000000000000005000000000101201dfc958f68d2ae3d6c66a5152ae104dabc2d81bd455c15f2f792b90f5dac4d8a74c9a5a495392bb04b71e8f3fc45180cdeb0d1f21216420390fffa0e7aa58de1207e5c0433da87b23db2ecc813a7f04a08ef061fc2fc67325d4a6872289c9138491a4724a688587ec5bdd98ea1f25f94d6cb697cd5386cd09f2136398b9a5290ce2014a74f7a2d8acd6887fe0306863529bc702cbec1ecf21f9f848bf57f07a5f78e57eba311bf2840dcc7734e049dd84516472a8805927d7172fe61f3342b580175",
			},
			wanted: want{
				alertKey: &Alert{
					Version:   0x01,
					Timestamp: 0x02,
					AlertType: MessageTypeBanPeer,
				},
				signatures: []string{
					"2014a74f7a2d8acd6887fe0306863529bc702cbec1ecf21f9f848bf57f07a5f78e57eba311bf2240dcc7734e049dd84516472a8805927d7172fe61f3342b580175",
					"2014a74f7a2d8acd6887fe0306863529bc702cbec1ecf21f9f848bf57f07a5f78e57eba311bf2240dcc7734e049dd84516472a881111117172fe61f3342b580175",
					"207e5c0433da87b23db2ecc813a7f04a08ef061fc2fc67325d4a6872289c9138491a4724a688587ec5bdd98ea1f25f94d6cb697cd5386cd09f2136398b9a5290ce",
				},
				data: "01000000123456789012345678901234567890121234567890123456789012345678901202000000050000000c093132372e302e302e310101",
			},
			wantErr: true,
		},
		{
			name: "valid invalidateblock alert",
			args: args{
				ak: "0100000001000000110000000000000007000000000000000000000006256d2b1da92365f44d905008b584c4729eb2018a4b67a2010120d15316302989ca03c6c8c42c4c42694e657a85652f4af4372c39e822503266f11a5342779779440e4e4bb61060ddd28f9a1c6e59239ad2e8eae8e561c95f5d22209f76606d6fa664184e5c0d0d53d1149e0c67a1d44757dc52fa2892abf61836dd45791c51f88904ce256272057b331b7995d227230c799bb9ecbf8d55b8da6fb520a40a9ffa90a08c2cabf54e998ac92102a31198b854f1d2acb677921a06c1b75e3c6670e57ea08d9a03abb7ddf086a8e61cc24bcb3cae777d14d1f9f8c297612b",
			},
			wanted: want{
				alertKey: &Alert{
					Version:   0x01,
					Timestamp: 0x11,
					AlertType: MessageTypeInvalidateBlock,
				},
				signatures: []string{
					"20d15316302989ca03c6c8c42c4c42694e657a85652f4af4372c39e822503266f11a5342779779440e4e4bb61060ddd28f9a1c6e59239ad2e8eae8e561c95f5d22",
					"209f76606d6fa664184e5c0d0d53d1149e0c67a1d44757dc52fa2892abf61836dd45791c51f88904ce256272057b331b7995d227230c799bb9ecbf8d55b8da6fb5",
					"20a40a9ffa90a08c2cabf54e998ac92102a31198b854f1d2acb677921a06c1b75e3c6670e57ea08d9a03abb7ddf086a8e61cc24bcb3cae777d14d1f9f8c297612b",
				},
				data: "0100000001000000110000000000000007000000000000000000000006256d2b1da92365f44d905008b584c4729eb2018a4b67a20101",
			},
		},
		{
			name: "valid unban alert",
			args: args{
				ak: "01000000010000002b00000000000000060000000b3132372e302e302e312f3001011fc03cd71dcfb817fac7000a5c20884c50f9c20ce353f428c531ecc8101ce6eae807e9b81a9ff09fa3e994451c8fb7b85d208daa90053fade9589eb680f683673f1f5902cc2c52877bf73aab98257cf310e4c08473a46261c9fb0a2df03de0d863b632cad667d97a4c20d5cc5081a55c9a6eca13a668d682e071af6f23d160a39d5120b78e14d2cde2feaf13674e0149d22ceaa2e6356dfb5790416901c7504e9483f64b3d812089c295f1d4a24358a894f1a8278aa0f4f3dc5141dd6a4c77a968be67",
			},
			wanted: want{
				alertKey: &Alert{
					Version:   0x01,
					Timestamp: 0x2b,
					AlertType: MessageTypeUnbanPeer,
				},
				signatures: []string{
					"1fc03cd71dcfb817fac7000a5c20884c50f9c20ce353f428c531ecc8101ce6eae807e9b81a9ff09fa3e994451c8fb7b85d208daa90053fade9589eb680f683673f",
					"1f5902cc2c52877bf73aab98257cf310e4c08473a46261c9fb0a2df03de0d863b632cad667d97a4c20d5cc5081a55c9a6eca13a668d682e071af6f23d160a39d51",
					"20b78e14d2cde2feaf13674e0149d22ceaa2e6356dfb5790416901c7504e9483f64b3d812089c295f1d4a24358a894f1a8278aa0f4f3dc5141dd6a4c77a968be67",
				},
				data: "01000000010000002b00000000000000060000000b3132372e302e302e312f300101",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//tt.wanted.alertKey.Signatures, _ = hex.DecodeString(tt.wanted.signature)
			tt.wanted.alertKey.Data, _ = hex.DecodeString(tt.wanted.data)

			ak, err := hex.DecodeString(tt.args.ak)
			if err != nil {
				t.Errorf("NewAlertKeyFromBytes() error decoding test hex string = %v", err)
				return
			}
			var got *Alert
			if got, err = NewAlertFromBytes(ak, nil); err != nil {
				t.Errorf("NewAlertKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got.Data, tt.wanted.alertKey.Data) && !tt.wantErr {
				t.Errorf("NewAlertKeyFromBytes().Data got = %#v, want %#v", got.Data, tt.wanted.alertKey.Data)
			}
			for i, s := range got.Signatures {
				sigHex := hex.EncodeToString(s)
				if sigHex != tt.wanted.signatures[i] && !tt.wantErr {
					t.Errorf("Signatures do not match got = %x, want %x", sigHex, tt.wanted.signatures[i])
				}

			}
			if got.AlertType != tt.wanted.alertKey.AlertType && !tt.wantErr {
				t.Errorf("NewAlertKeyFromBytes().AlertType got = %#v, want %#v", got.AlertType, tt.wanted.alertKey.AlertType)
			}
		})
	}
}*/
