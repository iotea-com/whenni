-- SCHEMA
CREATE DATABASE IF NOT EXISTS runtime_metrics;

-- TABLES
CREATE TABLE IF NOT EXISTS runtime_metrics.resources (
    -- Fields (core metrics) with efficient codecs
    `Timestamp` DateTime CODEC(Delta(8), ZSTD(1)),
    `CpuMilliCores` UInt64 CODEC(Delta(8), ZSTD(1)),
    `MemoryBytes` UInt64 CODEC(Delta(8), ZSTD(1)),
    `NetworkRxSpeed` UInt64 CODEC(Delta(8), ZSTD(1)),
    `NetworkTxSpeed` UInt64 CODEC(Delta(8), ZSTD(1)),
    `NetworkRxBytes` UInt64 CODEC(Delta(8), ZSTD(1)),
    `NetworkTxBytes` UInt64 CODEC(Delta(8), ZSTD(1)),
    
    -- Flattened attributes for simpler querying
    `ChannelId` String CODEC(ZSTD(1)),
    `ServiceName` String CODEC(ZSTD(1)),
    `Environment` String CODEC(ZSTD(1)),
    `ContainerId` String CODEC(ZSTD(1)),
    `ContainerName` String CODEC(ZSTD(1)),
    `ContainerImage` String CODEC(ZSTD(1)),
    `ContainerTag` String CODEC(ZSTD(1)),
    `PodName` String CODEC(ZSTD(1)),
    `Namespace` String CODEC(ZSTD(1))
) ENGINE = MergeTree()
PARTITION BY toDate(Timestamp)
ORDER BY (ChannelId, ServiceName, Timestamp)
TTL Timestamp + INTERVAL 7 DAY
SETTINGS 
    index_granularity = 8192,
    ttl_only_drop_parts = 1;