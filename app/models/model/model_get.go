package model

import (
	"context"
	"errors"
	"time"

	"github.com/mrz1836/go-datastore"
)

// Get will retrieve a model from the Datastore using the provided conditions
func Get(
	ctx context.Context,
	model BaseInterface,
	conditions map[string]interface{},
	timeout time.Duration,
	forceWriteDB bool,
) error {

	if timeout == 0 {
		timeout = DefaultDatabaseReadTimeout
	}

	// Attempt to Get the model (by model fields and given conditions)
	return model.Datastore().GetModel(ctx, model, conditions, timeout, forceWriteDB)
}

// GetModels will retrieve model(s) from the Datastore using the provided conditions
func GetModels(
	ctx context.Context,
	datastore datastore.ClientInterface,
	models interface{},
	conditions map[string]interface{},
	queryParams *datastore.QueryParams,
	timeout time.Duration,
) error {
	// Attempt to Get the model (by model fields and given conditions)
	return datastore.GetModels(ctx, models, conditions, queryParams, nil, timeout)
}

/*
// GetModelsAggregate will retrieve a count of the model(s) from the Datastore using the provided conditions
func GetModelsAggregate(
	ctx context.Context,
	datastore datastore.ClientInterface,
	models interface{},
	conditions map[string]interface{},
	aggregateColumn string,
	timeout time.Duration,
) (map[string]interface{}, error) {
	// Attempt to Get the model (by model fields & given conditions)
	return datastore.GetModelsAggregate(ctx, models, conditions, aggregateColumn, timeout)
}
*/

/*
// GetModelCount will retrieve a count of the model from the Datastore using the provided conditions
func GetModelCount(
	ctx context.Context,
	datastore datastore.ClientInterface,
	model interface{},
	conditions map[string]interface{},
	timeout time.Duration, //nolint:nolintlint,unparam // default timeout is passed most of the time
) (int64, error) {
	// Attempt to Get the model (by model fields & given conditions)
	return datastore.GetModelCount(ctx, model, conditions, timeout)
}
*/

// GetModelsByConditions will get models by given conditions
func GetModelsByConditions(ctx context.Context, modelName Name, modelItems interface{},
	metadata *Metadata, conditions *map[string]interface{}, queryParams *datastore.QueryParams,
	opts ...Options) error {

	dbConditions := map[string]interface{}{}

	if metadata != nil {
		dbConditions[MetadataField] = metadata
	}

	if conditions != nil && len(*conditions) > 0 {
		and := make([]map[string]interface{}, 0)
		if _, ok := dbConditions["$and"]; ok {
			and = dbConditions["$and"].([]map[string]interface{})
		}
		and = append(and, *conditions)
		dbConditions["$and"] = and
	}

	// Get the records
	if err := GetModels(
		ctx, NewBaseModel(modelName, opts...).Datastore(),
		modelItems, dbConditions, queryParams, DefaultDatabaseReadTimeout,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil
		}
		return err
	}

	return nil
}

/*
// getModelsAggregateByConditions will get aggregates of models by given conditions
func getModelsAggregateByConditions(ctx context.Context, modelName Name, models interface{},
	metadata *Metadata, conditions *map[string]interface{}, aggregateColumn string,
	opts ...Options) (map[string]interface{}, error) {

	dbConditions := map[string]interface{}{}

	if metadata != nil {
		dbConditions[MetadataField] = metadata
	}

	if conditions != nil && len(*conditions) > 0 {
		and := make([]map[string]interface{}, 0)
		if _, ok := dbConditions["$and"]; ok {
			and = dbConditions["$and"].([]map[string]interface{})
		}
		and = append(and, *conditions)
		dbConditions["$and"] = and
	}

	// Get the records
	results, err := GetModelsAggregate(
		ctx, NewBaseModel(modelName, opts...).Datastore(),
		models, dbConditions, aggregateColumn, DefaultDatabaseReadTimeout,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return results, nil
}
*/

/*
// getModelCountByConditions will get model counts (sums) from given conditions
func getModelCountByConditions(ctx context.Context, modelName Name, model interface{},
	metadata *Metadata, conditions *map[string]interface{}, opts ...Options) (int64, error) {

	dbConditions := map[string]interface{}{}

	if metadata != nil {
		dbConditions[MetadataField] = metadata
	}

	if conditions != nil && len(*conditions) > 0 {
		and := make([]map[string]interface{}, 0)
		if _, ok := dbConditions["$and"]; ok {
			and = dbConditions["$and"].([]map[string]interface{})
		}
		and = append(and, *conditions)
		dbConditions["$and"] = and
	}

	// Get the records
	count, err := GetModelCount(
		ctx, NewBaseModel(modelName, opts...).Datastore(),
		model, dbConditions, DefaultDatabaseReadTimeout,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}
*/
