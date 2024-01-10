-- Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

-- ../db/clickhouse/common/grpc-server/calls.sql

CREATE TABLE IF NOT EXISTS web_client_grpc_server.calls

CREATE TABLE IF NOT EXISTS web_client_grpc_server.call_queue
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'web_client.grpc_server.calls',
    kafka_group_name = 'web_client_grpc_server_calls_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/net/grpc/server/call_info:CallInfo';

CREATE MATERIALIZED VIEW IF NOT EXISTS web_client_grpc_server.call_consumer TO web_client_grpc_server.calls AS
SELECT
    ...
FROM web_client_grpc_server.call_queue;

-------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS web_client_grpc_server.calls
(
    id UUID,
    app Tuple(id UInt64, group_id UInt64, version LowCardinality(String), env LowCardinality(String)),
    app_session_id UInt64,
    grpc_server_id UInt16,
    status Enum8('unspecified' = 0, 'new' = 1, 'in_progress' = 2, 'success' = 3, 'failure' = 4),
    start_time DateTime64(6, 'UTC'),
    end_time Nullable(DateTime64(6, 'UTC')),
    elapsed_time_us Nullable(Int64), -- in microseconds
    full_method LowCardinality(String),
    content_type Array(LowCardinality(String)),
    user_agent Array(LowCardinality(String)),
    is_operation_successful Nullable(Bool),
    status_code Nullable(UInt32),
    _version_stamp UInt8,

    INDEX id_idx id TYPE bloom_filter GRANULARITY 1,
    INDEX app_id_idx tupleElement(app, 'id') TYPE set(0) GRANULARITY 1,
    INDEX app_group_id_idx tupleElement(app, 'group_id') TYPE set(0) GRANULARITY 1,
    INDEX app_session_id_idx app_session_id TYPE set(0) GRANULARITY 1,
    INDEX grpc_server_id_idx grpc_server_id TYPE set(0) GRANULARITY 1,
    INDEX status_idx status TYPE set(0) GRANULARITY 1,
    INDEX elapsed_time_us_idx elapsed_time_us TYPE minmax GRANULARITY 1,
    INDEX full_method_idx lowerUTF8(full_method) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX content_type_idx arrayMap(x -> lowerUTF8(x), content_type) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX user_agent_idx arrayMap(x -> lowerUTF8(x), user_agent) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX is_operation_successful_idx is_operation_successful TYPE set(0) GRANULARITY 1,
    INDEX status_code_idx status_code TYPE set(0) GRANULARITY 1,
    INDEX _version_stamp_idx _version_stamp TYPE set(0) GRANULARITY 1
)
ENGINE = ReplacingMergeTree(_version_stamp)
PARTITION BY toYYYYMM(start_time)
ORDER BY (start_time, id)
SETTINGS index_granularity=8192;

CREATE TABLE IF NOT EXISTS web_client_grpc_server.call_queue
(
    id String,
    app Tuple(id UInt64, group_id UInt64, version String, env String),
    app_session_id UInt64,
    grpc_server_id UInt16,
    status Enum8('unspecified' = 0, 'new' = 1, 'in_progress' = 2, 'success' = 3, 'failure' = 4),
    start_time Int64, -- in microseconds
    end_time Nullable(Int64), -- in microseconds
    elapsed_time_us Nullable(Int64), -- in microseconds
    full_method String,
    content_type Array(String),
    user_agent Array(String),
    is_operation_successful Nullable(Bool),
    status_code Nullable(UInt32)
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'web_client.grpc_server.calls',
    kafka_group_name = 'web_client_grpc_server_calls_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/net/grpc/server/call_info:CallInfo';

CREATE MATERIALIZED VIEW IF NOT EXISTS web_client_grpc_server.call_consumer TO web_client_grpc_server.calls AS
SELECT
    toUUID(id) AS id,
    app,
    app_session_id,
    grpc_server_id,
    status,
    fromUnixTimestamp64Micro(start_time, 'UTC') AS start_time,
    fromUnixTimestamp64Micro(end_time, 'UTC') AS end_time,
    elapsed_time_us,
    full_method,
    content_type,
    user_agent,
    is_operation_successful,
    status_code,
    if(end_time IS NOT NULL, 2, 1) AS _version_stamp
FROM web_client_grpc_server.call_queue;
