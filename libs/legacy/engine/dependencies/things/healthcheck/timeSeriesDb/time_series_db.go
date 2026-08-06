package timeSeriesDbHealthcheck

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"

	"fmt"
)

func InfluxDbDatabase(attrs *things.InfluxDbDatabase) error {
	// Connect to InfluxDB
	databaseUrl := fmt.Sprintf("%s://%s:%d", attrs.Protocol, attrs.Host, attrs.Port)
	client := influxdb2.NewClient(databaseUrl, attrs.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Ping the server
	ok, err := client.Ping(ctx)
	if !ok || err != nil {
		return fmt.Errorf("failed to ping the InfluxDB server at %s:%d: %s", attrs.Host, attrs.Port, err)
	}

	return nil
}

func ClickhouseDatabase(attrs *things.ClickhouseDatabase) error {
	// Connect to ClickHouse
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr:     []string{fmt.Sprintf("%s:%d", attrs.Host, attrs.Port)},
		Protocol: clickhouse.Native,
		Auth: clickhouse.Auth{
			Database: attrs.Database,
			Username: attrs.Username,
			Password: attrs.Password,
		},
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	})
	if err != nil {
		return fmt.Errorf("could not connect to ClickHouse: %s", err)
	}

	// Give it some time to connect
	pingTimeout := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	// Ping the server
	err = conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping the ClickHouse server at %s:%d: %s", attrs.Host, attrs.Port, err)
	}

	return nil
}
