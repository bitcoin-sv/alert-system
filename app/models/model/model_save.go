package model

import (
	"context"
	"fmt"

	"github.com/mrz1836/go-datastore"
	"github.com/pkg/errors"
)

// Save will save the model(s) into the Datastore
func Save(ctx context.Context, model BaseInterface) (err error) {

	// Check for datastore
	ds := model.Datastore()
	if ds == nil {
		return ErrMissingDatastore
	}

	// Create new Datastore transaction
	// NOTE: we need this to be in a callback context for Mongo
	return ds.NewTx(ctx, func(tx *datastore.Transaction) (err error) {

		// Fire the before hooks (parent model)
		if model.IsNew() {
			if err = model.BeforeCreating(ctx); err != nil {
				return
			}
		} else {
			if err = model.BeforeUpdating(ctx); err != nil {
				return
			}
		}

		// Set the record's timestamps
		model.SetRecordTime(model.IsNew())

		// Sync the list of models to Save
		modelsToSave := append(make([]BaseInterface, 0), model)

		// Add any child models (fire before hooks)
		if children := model.ChildModels(); len(children) > 0 {
			for _, child := range children {
				if child.IsNew() {
					if err = child.BeforeCreating(ctx); err != nil {
						return
					}
				} else {
					if err = child.BeforeUpdating(ctx); err != nil {
						return
					}
				}

				// Set the record's timestamps
				child.SetRecordTime(child.IsNew())
			}

			// Add to list for saving
			modelsToSave = append(modelsToSave, children...)
		}

		// Logs for saving models
		model.DebugLog(ctx, fmt.Sprintf("saving %d models...", len(modelsToSave)))

		// Save all models (or fail!)
		for index := range modelsToSave {
			// modelsToSave[index].DebugLog(ctx, fmt.Sprintf("starting to save model: %s id: %d", modelsToSave[index].Name(), modelsToSave[index].GetID()))
			if err = modelsToSave[index].Datastore().SaveModel(
				ctx, modelsToSave[index], tx, modelsToSave[index].IsNew(), false,
			); err != nil {
				return
			}
		}

		// Commit all the model(s) if needed
		if tx.CanCommit() {
			// model.DebugLog(ctx, "committing db transaction...")
			if err = tx.Commit(); err != nil {
				return
			}
		}

		// Fire after hooks (only on commit success)
		var afterErr error
		for index := range modelsToSave {
			if modelsToSave[index].IsNew() {
				modelsToSave[index].NotNew() // NOTE: calling it before this method... after created assumes it's been saved already
				afterErr = modelsToSave[index].AfterCreated(ctx)
			} else {
				afterErr = modelsToSave[index].AfterUpdated(ctx)
			}
			if afterErr != nil {
				if err == nil { // First error - set the error
					err = afterErr
				} else { // Got more than one error, wrap it!
					err = errors.Wrap(err, afterErr.Error())
				}
			}
			// modelToSave.NotNew() // NOTE: moved to above from here
		}

		return
	})
}

// BeginSaveWithTx will start saving the model(s) into the Datastore
func BeginSaveWithTx(ctx context.Context, tx *datastore.Transaction, model BaseInterface) (modelsToSave []BaseInterface, err error) {

	// Check for datastore
	ds := model.Datastore()
	if ds == nil {
		return nil, ErrMissingDatastore
	}

	// Fire the before hooks (parent model)
	if model.IsNew() {
		if err = model.BeforeCreating(ctx); err != nil {
			return
		}
	} else {
		if err = model.BeforeUpdating(ctx); err != nil {
			return
		}
	}

	// Set the record's timestamps
	model.SetRecordTime(model.IsNew())

	// Sync the list of models to Save
	modelsToSave = append(make([]BaseInterface, 0), model)

	// Add any child models (fire before hooks)
	if children := model.ChildModels(); len(children) > 0 {
		for _, child := range children {
			if child.IsNew() {
				if err = child.BeforeCreating(ctx); err != nil {
					return
				}
			} else {
				if err = child.BeforeUpdating(ctx); err != nil {
					return
				}
			}

			// Set the record's timestamps
			child.SetRecordTime(child.IsNew())
		}

		// Add to list for saving
		modelsToSave = append(modelsToSave, children...)
	}

	// Logs for saving models
	// model.DebugLog(ctx, fmt.Sprintf("saving %d models...", len(modelsToSave)))

	// Save all models (or fail!)
	for index := range modelsToSave {
		// modelsToSave[index].DebugLog(ctx, fmt.Sprintf("starting to save model: %s id: %d", modelsToSave[index].Name(), modelsToSave[index].GetID()))
		if err = modelsToSave[index].Datastore().SaveModel(
			ctx, modelsToSave[index], tx, modelsToSave[index].IsNew(), false,
		); err != nil {
			return
		}
	}

	return
}

// CompleteSaveWithTx will finish saving the model(s) into the Datastore
func CompleteSaveWithTx(ctx context.Context, tx *datastore.Transaction, modelsToSave []BaseInterface) (err error) {

	// Commit all the model(s) if needed
	if tx.CanCommit() {
		// modelsToSave[0].DebugLog(ctx, "committing db transaction...")
		if err = tx.Commit(); err != nil {
			return
		}
	}

	// Fire after hooks (only on commit success)
	var afterErr error
	for index := range modelsToSave {
		if modelsToSave[index].IsNew() {
			modelsToSave[index].NotNew() // NOTE: calling it before this method... after created assumes it's been saved already
			afterErr = modelsToSave[index].AfterCreated(ctx)
		} else {
			afterErr = modelsToSave[index].AfterUpdated(ctx)
		}
		if afterErr != nil {
			if err == nil { // First error - set the error
				err = afterErr
			} else { // Got more than one error, wrap it!
				err = errors.Wrap(err, afterErr.Error())
			}
		}
		// modelToSave.NotNew() // NOTE: moved to above from here
	}

	return
}
