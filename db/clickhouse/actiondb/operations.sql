-- Copyright 2023 Alexey Lavrenchenko. All rights reserved.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
-- 	http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

CREATE TABLE IF NOT EXISTS actiondb.operations
(
    id UUID,
    app Tuple(id UInt64, group_id UInt64, version LowCardinality(String), env LowCardinality(String)),
    app_session_id UInt64,
    tran_id UUID,
    action_id UUID,
    `type` UInt64,
    category Enum8('unspecified' = 0, 'common' = 1, 'identity' = 2, 'database' = 3, 'cache_storage' = 4),
    `group` UInt64,
    parent_operation_id Nullable(UUID),
    params Nullable(String),
    param_map Map(LowCardinality(String), String),
    created_at DateTime64(6, 'UTC'),
    status Enum8('unspecified' = 0, 'new' = 1, 'in_progress' = 2, 'success' = 3, 'failure' = 4),
    start_time DateTime64(6, 'UTC'),
    end_time Nullable(DateTime64(6, 'UTC')),
    elapsed_time_us Nullable(Int64), -- in microseconds
    _version_stamp UInt8,

    INDEX id_idx id TYPE bloom_filter GRANULARITY 1,
    INDEX app_id_idx tupleElement(app, 'id') TYPE set(0) GRANULARITY 1,
    INDEX app_group_id_idx tupleElement(app, 'group_id') TYPE set(0) GRANULARITY 1,
    INDEX app_session_id_idx app_session_id TYPE set(0) GRANULARITY 1,
    INDEX tran_id_idx tran_id TYPE bloom_filter GRANULARITY 1,
    INDEX action_id_idx action_id TYPE bloom_filter GRANULARITY 1,
    INDEX type_idx `type` TYPE set(0) GRANULARITY 1,
    INDEX category_idx category TYPE set(0) GRANULARITY 1,
    INDEX group_idx `group` TYPE set(0) GRANULARITY 1,
    INDEX parent_operation_id_idx parent_operation_id TYPE bloom_filter GRANULARITY 1,
    INDEX param_map_key_idx mapKeys(param_map) TYPE set(0) GRANULARITY 1,
    INDEX param_map_value_idx mapValues(param_map) TYPE bloom_filter GRANULARITY 1,
    INDEX status_idx status TYPE set(0) GRANULARITY 1,
    INDEX elapsed_time_us_idx elapsed_time_us TYPE minmax GRANULARITY 1,
    INDEX _version_stamp_idx _version_stamp TYPE set(0) GRANULARITY 1
)
ENGINE = ReplacingMergeTree(_version_stamp)
PARTITION BY toYYYYMM(created_at)
ORDER BY (created_at, id)
SETTINGS index_granularity=8192;

CREATE TABLE IF NOT EXISTS actiondb.operation_queue
(
    id String,
    app Tuple(id UInt64, group_id UInt64, version String, env String),
    app_session_id UInt64,
    tran_id String,
    action_id String,
    `type` UInt64,
    category Enum8('unspecified' = 0, 'common' = 1, 'identity' = 2, 'database' = 3, 'cache_storage' = 4),
    `group` UInt64,
    parent_operation_id Nullable(String),
    params Nullable(String),
    created_at Int64, -- in microseconds
    status Enum8('unspecified' = 0, 'new' = 1, 'in_progress' = 2, 'success' = 3, 'failure' = 4),
    start_time Int64, -- in microseconds
    end_time Nullable(Int64), -- in microseconds
    elapsed_time_us Nullable(Int64) -- in microseconds
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'base.operations',
    kafka_group_name = 'base_operations_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/actions/operation:Operation';

CREATE MATERIALIZED VIEW IF NOT EXISTS actiondb.operation_consumer TO actiondb.operations AS
SELECT
    toUUID(id) AS id,
    app,
    app_session_id,
    toUUID(tran_id) AS tran_id,
    toUUID(action_id) AS action_id,
    type,
    category,
    group,
    toUUID(parent_operation_id) AS parent_operation_id,
    params,
    JSONExtract(ifNull(params, ''), 'Map(LowCardinality(String), String)') AS param_map,
    fromUnixTimestamp64Micro(created_at, 'UTC') AS created_at,
    status,
    fromUnixTimestamp64Micro(start_time, 'UTC') AS start_time,
    fromUnixTimestamp64Micro(end_time, 'UTC') AS end_time,
    elapsed_time_us,
    if(end_time IS NOT NULL, 2, 1) AS _version_stamp
FROM actiondb.operation_queue;
