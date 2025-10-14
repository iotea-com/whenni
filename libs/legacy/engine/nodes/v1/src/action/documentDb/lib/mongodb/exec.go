package mongodbActionNode

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/iotea-com/iotea/libs/engine/template"

	"go.mongodb.org/mongo-driver/bson"
)

func (n *MongoDbActionSubnode) Exec(params node.ExecParams) node.Error {
	// Client library works with a registered callback. In this loop we simply
	// await for a context cancellation, or an error
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			invalid, err := n.handle(ioData)
			if err != nil {
				if invalid {
					// Safely notify the runtime that data was processed
					n.notifyChannel <- node.Notification{
						Type:   node.NotifyDataInvalid,
						Reason: err.Error(),
					}
				} else {
					return node.Error{
						Type:   node.FatalError,
						Reason: err.Error(),
					}
				}
			} else {
				// Safely notify the runtime that data was processed
				n.notifyChannel <- node.Notification{
					Type: node.NotifyDataProcessed,
				}
			}
		}
	}

	return node.Error{
		Type: node.NoError,
	}
}

func (n *MongoDbActionSubnode) handle(ioData node.IoData) (bool, error) {
	// Convert IoData to map for filter and document
	dataMap, err := template.ConvertIoDataToMap(ioData)
	if err != nil {
		n.logger.Error().Ctx(ioData.Ctx).Msgf("Failed to convert data to map: %v", err)
		return true, fmt.Errorf("failed to convert data to map: %v", err) // Validation error
	}

	filter := (func() bson.M {
		if n.filterTemplate.Value == "" {
			return nil
		}

		// Preprocess the filter template
		n.filterTemplate.Preprocess()

		// Fill the filter template with data from IoData
		filterString, err := n.filterTemplate.Fill(dataMap)
		if err != nil {
			n.logger.Error().Ctx(ioData.Ctx).Msgf("Failed to fill the filter template: %v", err)
			return nil
		}

		// Parse the filter string into a BSON object
		var b bson.M
		if err := bson.UnmarshalExtJSON([]byte(*filterString), false, &b); err != nil {
			n.logger.Error().Ctx(ioData.Ctx).Msgf("Could not convert filter into BSON: %s", err)
			return nil
		}

		return b
	})()

	document := (func() bson.M {
		if n.documentTemplate.Value == "" {
			return nil
		}

		// Preprocess the document template
		n.documentTemplate.Preprocess()

		// Fill the document template with data from IoData
		documentString, err := n.documentTemplate.Fill(dataMap)
		if err != nil {
			n.logger.Error().Ctx(ioData.Ctx).Msgf("Failed to fill the document template: %v", err)
			return nil
		}

		// Parse the filter string into a BSON object
		var b bson.M
		if err := bson.UnmarshalExtJSON([]byte(*documentString), false, &b); err != nil {
			n.logger.Error().Ctx(ioData.Ctx).Msgf("Could not convert document into BSON: %s", err)
			return nil
		}

		return b
	})()

	n.logger.Info().Ctx(ioData.Ctx).Msgf("Readying MQL query to server at %s:%d", n.config.MongoDb.Host, n.config.MongoDb.Port)

	// Ping the MongoDB server
	dbCtx := context.Background()
	defer dbCtx.Done()

	err = n.mongoClient.Ping(dbCtx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to ping the MongoDB server at %s:%d: %s", n.config.MongoDb.Host, n.config.MongoDb.Port, err) // Fatal error
	}

	// Execute the parsed BSON query
	collection := n.mongoClient.Database(n.config.Database).Collection(n.config.Collection)
	results, err := (func() ([]bson.M, error) {
		switch n.config.Method {
		case string(MongoDbQueryMethodFind):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the find")
			}

			cursor, err := collection.Find(dbCtx, filter)
			if err != nil {
				return nil, err
			}
			defer cursor.Close(dbCtx)

			results := []bson.M{}
			for cursor.Next(dbCtx) {
				var result bson.M
				if err := cursor.Decode(&result); err != nil {
					n.logger.Error().Ctx(ioData.Ctx).Msgf("Failed to decode the results: %s", err)
					continue
				}

				results = append(results, result)
			}
			return results, nil
		case string(MongoDbQueryMethodFindOne):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the find")
			}

			singleResult := collection.FindOne(dbCtx, filter)
			var result bson.M
			if err := singleResult.Decode(&result); err != nil {
				n.logger.Error().Ctx(ioData.Ctx).Msgf("Failed to decode the result: %s", err)
				return nil, err
			}
			return []bson.M{result}, nil
		case string(MongoDbQueryMethodFindOneAndDelete):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the find")
			}

			singleResult := collection.FindOneAndDelete(dbCtx, filter)
			var result bson.M
			if err := singleResult.Decode(&result); err != nil {
				n.logger.Error().Ctx(ioData.Ctx).Msgf("Failed to decode the result: %s", err)
				return nil, err
			}
			return []bson.M{result}, nil
		case string(MongoDbQueryMethodFindOneAndReplace):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the find")
			}
			if document == nil {
				return nil, fmt.Errorf("no document was provided for the replace")
			}

			singleResult := collection.FindOneAndReplace(dbCtx, filter, document)
			var result bson.M
			if err := singleResult.Decode(&result); err != nil {
				n.logger.Error().Ctx(ioData.Ctx).Msgf("Failed to decode the result: %s", err)
				return nil, err
			}
			return []bson.M{result}, nil
		case string(MongoDbQueryMethodFindOneAndUpdate):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the find")
			}
			if document == nil {
				return nil, fmt.Errorf("no document was provided for the update")
			}

			singleResult := collection.FindOneAndUpdate(dbCtx, filter, document)
			var result bson.M
			if err := singleResult.Decode(&result); err != nil {
				n.logger.Error().Ctx(ioData.Ctx).Msgf("Failed to decode the result: %s", err)
				return nil, err
			}
			return []bson.M{result}, nil
		case string(MongoDbQueryMethodDeleteMany):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the find")
			}

			_, err := collection.DeleteMany(dbCtx, filter)
			if err != nil {
				return nil, err
			}
			return nil, nil
		case string(MongoDbQueryMethodDeleteOne):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the insert")
			}

			_, err := collection.DeleteOne(dbCtx, filter)
			if err != nil {
				return nil, err
			}
			return nil, nil
		case string(MongoDbQueryMethodInsertMany):
			return nil, fmt.Errorf("method InsertMany is not supported at this time")
		case string(MongoDbQueryMethodInsertOne):
			if document == nil {
				return nil, fmt.Errorf("no document was provided for the insert")
			}

			_, err := collection.InsertOne(dbCtx, document)
			return nil, err
		case string(MongoDbQueryMethodReplaceOne):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the find")
			}
			if document == nil {
				return nil, fmt.Errorf("no document was provided for the replace")
			}

			_, err := collection.ReplaceOne(dbCtx, filter, document)
			return nil, err
		case string(MongoDbQueryMethodUpdateMany):
			return nil, fmt.Errorf("method UpdateMany is not supported at this time")
		case string(MongoDbQueryMethodUpdateOne):
			if filter == nil {
				return nil, fmt.Errorf("no filter was provided for the find")
			}
			if document == nil {
				return nil, fmt.Errorf("no document was provided for the update")
			}

			_, err := collection.UpdateOne(dbCtx, filter, document)

			return nil, err
		default:
			return nil, fmt.Errorf("invalid method '%s'", n.config.Method)
		}
	})()

	if err != nil {
		return false, fmt.Errorf("failed to execute query: %v", err) // Fatal error
	}

	n.logger.Info().Ctx(ioData.Ctx).Msgf("Successfully executed the %s query", n.config.Method)

	resultsBytes, err := json.Marshal(results)
	if err != nil {
		return false, fmt.Errorf("failed to marshal results: %v", err) // Fatal error
	}

	outputData := node.IoData{
		Data: resultsBytes,
		Type: node.BytesDataType,
		Ctx:  ioData.Ctx,
	}

	select {
	case n.outputChannels[0].Channel <- outputData:
	default:
		return false, nil
	}

	return false, nil
}
