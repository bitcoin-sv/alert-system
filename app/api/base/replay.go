package base

/*// ReplayResponse is the response for the replay endpoint
type ReplayResponse struct {
}

// ReplayEndpoint will replay all alerts from the database
func ReplayEndpoint(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	res := HealthResponse{
		Synced: true, // TODO actually fetch this state from the DB somehow, or from the server struct

	}
	go func() {
		_ = replayAlerts() // todo handle error
	}()
	_ = apirouter.ReturnJSONEncode(w, http.StatusAccepted, json.NewEncoder(w), res, []string{"alert", "synced", "sequence"})
}

// replayAlerts will replay all alerts from the database
func replayAlerts() error {
	ctx := context.Background()
	log := gocore.Log("alert-system-force-replay")

	db, err := database.NewSqliteClient(ctx, log)
	if err != nil {
		return err
	}

	var a *model.AlertMessage
	if a, err = db.GetLatestAlert(ctx); err != nil {
		return err
	}
	latestSequence := a.SequenceNumber
	for i := uint32(1); i <= latestSequence; i++ {
		if a, err = db.GetAlertBySequenceNumber(ctx, i); err != nil {
			return err
		}
		// Do the alert
		var raw []byte
		if raw, err = hex.DecodeString(a.Raw); err != nil {
			return err
		}

		var alertMessage *alert.Alert
		// todo this needs a valid config!
		if alertMessage, err = alert.NewAlertFromBytes(raw, nil); err != nil {
			// probably want to ban this peer
			return err
		}

		ak := alertMessage.ProcessAlertMessage()
		if err = ak.Read(alertMessage.AlertMessage); err != nil {
			return err
		}
		if err = ak.Do(ctx); err != nil {
			return err
		}
	}
	return nil
}*/
