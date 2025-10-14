package things

type ThingCategory string

func (tc ThingCategory) String() string {
	return string(tc)
}

func (tc *ThingCategory) StringOrNil() *string {
	if tc == nil {
		return nil
	}

	return (*string)(tc)
}

var ThingCategories = []ThingCategory{
	HttpServerThingCategory,
	MqttBrokerThingCategory,
	MqttClientThingCategory,
	KafkaClusterThingCategory,
	KafkaProducerThingCategory,
	KafkaConsumerThingCategory,
	NatsServerThingCategory,
	NatsClientThingCategory,
	InfluxDbDatabaseThingCategory,
	ClickhouseDatabaseThingCategory,
	AwsSNSThingCategory,
	AwsSESThingCategory,
	SendgridClientThingCategory,
}

func ParseCategoryString(category string) (*ThingCategory, bool) {
	for _, tc := range ThingCategories {
		if string(tc) == category {
			return &tc, true
		}
	}
	return nil, false
}
