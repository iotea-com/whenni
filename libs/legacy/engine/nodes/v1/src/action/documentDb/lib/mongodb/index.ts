export enum MongoDbQueryMethod {
	"Find" = "Find",
	"FindOne" = "FindOne",
	"FindOneAndUpdate" = "FindOneAndUpdate",
	"FindOneAndReplace" = "FindOneAndReplace",
	"FindOneAndDelete" = "FindOneAndDelete",
	"InsertOne" = "InsertOne",
	// "InsertMany" = "InsertMany",
	"UpdateOne" = "UpdateOne",
	// "UpdateMany" = "UpdateMany",
	"ReplaceOne" = "ReplaceOne",
	"DeleteOne" = "DeleteOne",
	"DeleteMany" = "DeleteMany",
}

export const MongoDbQueryMethodOptions = [
  {
    label: MongoDbQueryMethod.Find,
    value: MongoDbQueryMethod.Find,
  },
  {
    label: MongoDbQueryMethod.FindOne,
    value: MongoDbQueryMethod.FindOne,
  },
  {
    label: MongoDbQueryMethod.FindOneAndUpdate,
    value: MongoDbQueryMethod.FindOneAndUpdate,
  },
  {
    label: MongoDbQueryMethod.FindOneAndReplace,
    value: MongoDbQueryMethod.FindOneAndReplace,
  },
  {
    label: MongoDbQueryMethod.FindOneAndDelete,
    value: MongoDbQueryMethod.FindOneAndDelete,
  },
  {
    label: MongoDbQueryMethod.InsertOne,
    value: MongoDbQueryMethod.InsertOne,
  },
  // {
  //   label: MongoDbQueryMethod.InsertMany,
  //   value: MongoDbQueryMethod.InsertMany,
  // },
  {
    label: MongoDbQueryMethod.UpdateOne,
    value: MongoDbQueryMethod.UpdateOne,
  },
  // {
  //   label: MongoDbQueryMethod.UpdateMany,
  //   value: MongoDbQueryMethod.UpdateMany,
  // },
  {
    label: MongoDbQueryMethod.ReplaceOne,
    value: MongoDbQueryMethod.ReplaceOne,
  },
  {
    label: MongoDbQueryMethod.DeleteOne,
    value: MongoDbQueryMethod.DeleteOne,
  },
  {
    label: MongoDbQueryMethod.DeleteMany,
    value: MongoDbQueryMethod.DeleteMany,
  },
]

export type MongoDbActionSubnodeConfig = {
  "mongoDbServer::thing"?: string
  database?: string
  collection?: string
  queryMethod?: MongoDbQueryMethod
  "input::model"?: string
  filter?: string
  document?: string
  defaultValues?: Record<string, any>
}
