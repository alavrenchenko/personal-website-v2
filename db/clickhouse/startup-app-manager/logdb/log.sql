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

-- ../db/clickhouse/common/logdb/log.sql

CREATE TABLE IF NOT EXISTS startup_app_manager_logdb.log

CREATE TABLE IF NOT EXISTS startup_app_manager_logdb.log_queue
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'startup_app_manager.log',
    kafka_group_name = 'startup_app_manager_log_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/logging/log_entry:LogEntry';

CREATE MATERIALIZED VIEW IF NOT EXISTS startup_app_manager_logdb.log_consumer TO startup_app_manager_logdb.log AS
SELECT
    ...
FROM startup_app_manager_logdb.log_queue;

-------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS startup_app_manager_logdb.log
(
    id UUID,
    timestamp DateTime64(6, 'UTC'),
    app Tuple(
        id UInt64,
        group_id UInt64,
        version LowCardinality(String),
        env LowCardinality(String)
    ),
    agent Tuple(
        name LowCardinality(String),
        `type` LowCardinality(String),
        version LowCardinality(String)
    ),
    logging_session_id UInt64,
    app_session_id Nullable(UInt64),
    tran Tuple(
        id Nullable(UUID)
    ) DEFAULT tuple(NULL),
    action Tuple(
        id Nullable(UUID),
        `type` UInt64,
        category Enum8('unspecified' = 0, 'common' = 1, 'http' = 2, 'grpc' = 3),
        group UInt64
    ) DEFAULT tuple(NULL, 0, 0, 0),
    operation Tuple(
        id Nullable(UUID),
        `type` UInt64,
        category Enum8('unspecified' = 0, 'common' = 1, 'identity' = 2, 'database' = 3, 'cache_storage' = 4),
        group UInt64
    ) DEFAULT tuple(NULL, 0, 0, 0),
    level Enum8('trace' = 0, 'debug' = 1, 'info' = 2, 'warning' = 3, 'error' = 4, 'fatal' = 5),
    category LowCardinality(String),
    event Tuple(
        id UInt64,
        name LowCardinality(String),
        category Enum8('unknown' = 0, 'common' = 1, 'configuration' = 2, 'identity' = 3, 'database' = 4, 'cache_storage' = 5, 'network' = 6),
        group UInt64
    ),
    error Tuple(
        code UInt64,
        message String,
        `type` LowCardinality(String),
        category Enum8('unspecified' = 0, 'common' = 1, 'api' = 2, 'database' = 3),
        stack_trace String,
        original_error Tuple(
            code UInt64,
            message String,
            `type` LowCardinality(String),
            category Enum8('unspecified' = 0, 'common' = 1, 'api' = 2, 'database' = 3),
            stack_trace String
        )
    ) DEFAULT tuple(0, '', '', 0, '', (0, '', '', 0, '')), -- NoError = 0
    message String,
    fields Nullable(String),
    field_map Map(LowCardinality(String), String),

    INDEX id_idx id TYPE bloom_filter GRANULARITY 1,
    INDEX app_id_idx tupleElement(app, 'id') TYPE set(0) GRANULARITY 1,
    INDEX app_group_id_idx tupleElement(app, 'group_id') TYPE set(0) GRANULARITY 1,
    INDEX agent_name_idx tupleElement(agent, 'name') TYPE set(0) GRANULARITY 1,
    INDEX logging_session_id_idx logging_session_id TYPE set(0) GRANULARITY 1,
    INDEX app_session_id_idx app_session_id TYPE set(0) GRANULARITY 1,
    INDEX tran_id_idx tupleElement(tran, 'id') TYPE bloom_filter GRANULARITY 1,
    INDEX action_id_idx tupleElement(action, 'id') TYPE bloom_filter GRANULARITY 1,
    INDEX action_type_idx tupleElement(action, 'type') TYPE set(0) GRANULARITY 1,
    INDEX action_category_idx tupleElement(action, 'category') TYPE set(0) GRANULARITY 1,
    INDEX action_group_idx tupleElement(action, 'group') TYPE set(0) GRANULARITY 1,
    INDEX operation_id_idx tupleElement(operation, 'id') TYPE bloom_filter GRANULARITY 1,
    INDEX operation_type_idx tupleElement(operation, 'type') TYPE set(0) GRANULARITY 1,
    INDEX operation_category_idx tupleElement(operation, 'category') TYPE set(0) GRANULARITY 1,
    INDEX operation_group_idx tupleElement(operation, 'group') TYPE set(0) GRANULARITY 1,
    INDEX level_idx level TYPE set(0) GRANULARITY 1,
    INDEX category_idx lowerUTF8(category) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX event_id_idx tupleElement(event, 'id') TYPE set(0) GRANULARITY 1,
    INDEX event_name_idx tupleElement(event, 'name') TYPE set(0) GRANULARITY 1,
    INDEX event_category_idx tupleElement(event, 'category') TYPE set(0) GRANULARITY 1,
    INDEX event_group_idx tupleElement(event, 'group') TYPE set(0) GRANULARITY 1,
    INDEX error_code_idx tupleElement(error, 'code') TYPE set(0) GRANULARITY 1,
    INDEX error_message_idx lowerUTF8(tupleElement(error, 'message')) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX error_type_idx tupleElement(error, 'type') TYPE set(0) GRANULARITY 1,
    INDEX error_category_idx tupleElement(error, 'category') TYPE set(0) GRANULARITY 1,
    INDEX error_original_error_code_idx tupleElement(tupleElement(error, 'original_error'), 'code') TYPE set(0) GRANULARITY 1,
    INDEX error_original_error_message_idx lowerUTF8(tupleElement(tupleElement(error, 'original_error'), 'message')) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX error_original_error_type_idx tupleElement(tupleElement(error, 'original_error'), 'type') TYPE set(0) GRANULARITY 1,
    INDEX error_original_error_category_idx tupleElement(tupleElement(error, 'original_error'), 'category') TYPE set(0) GRANULARITY 1,
    INDEX message_idx lowerUTF8(message) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX field_map_key_idx mapKeys(field_map) TYPE set(0) GRANULARITY 1,
    INDEX field_map_value_idx mapValues(field_map) TYPE bloom_filter GRANULARITY 1
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, id)
SETTINGS index_granularity=8192;

CREATE TABLE IF NOT EXISTS startup_app_manager_logdb.log_queue
(
    id String,
    timestamp Int64, -- in microseconds
    app Tuple(id UInt64, group_id UInt64, version String, env String),
    agent Tuple(name String, `type` String, version String),
    logging_session_id UInt64,
    app_session_id Nullable(UInt64),
    tran Tuple(
        id Nullable(String)
    ),
    action Tuple(
        id Nullable(String),
        `type` UInt64,
        category Enum8('unspecified' = 0, 'common' = 1, 'http' = 2, 'grpc' = 3),
        group UInt64
    ),
    operation Tuple(
        id Nullable(String),
        `type` UInt64,
        category Enum8('unspecified' = 0, 'common' = 1, 'identity' = 2, 'database' = 3, 'cache_storage' = 4),
        group UInt64
    ),
    level Enum8('trace' = 0, 'debug' = 1, 'info' = 2, 'warning' = 3, 'error' = 4, 'fatal' = 5),
    category String,
    event Tuple(
        id UInt64,
        name String,
        category Enum8('unknown' = 0, 'common' = 1, 'configuration' = 2, 'identity' = 3, 'database' = 4, 'cache_storage' = 5, 'network' = 6),
        group UInt64
    ),
    error Tuple(
        code UInt64,
        message String,
        `type` String,
        category Enum8('unspecified' = 0, 'common' = 1, 'api' = 2, 'database' = 3),
        stack_trace String,
        original_error Tuple(
            code UInt64,
            message String,
            `type` String,
            category Enum8('unspecified' = 0, 'common' = 1, 'api' = 2, 'database' = 3),
            stack_trace String
        )
    ), -- NoError = 0
    message String,
    fields Nullable(String)
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'startup_app_manager.log',
    kafka_group_name = 'startup_app_manager_log_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/logging/log_entry:LogEntry';

CREATE MATERIALIZED VIEW IF NOT EXISTS startup_app_manager_logdb.log_consumer TO startup_app_manager_logdb.log AS
SELECT
    toUUID(id) AS id,
    fromUnixTimestamp64Micro(timestamp, 'UTC') AS timestamp,
    app,
    agent,
    logging_session_id,
    app_session_id,
    tuple(toUUID(tupleElement(tran, 'id'))) AS tran,
    if(tupleElement(action, 'id') IS NOT NULL, (toUUID(tupleElement(action, 'id')), tupleElement(action, 'type'), tupleElement(action, 'category'), tupleElement(action, 'group')), (NULL, 0, 0, 0)) AS action,
    if(tupleElement(operation, 'id') IS NOT NULL, (toUUID(tupleElement(operation, 'id')), tupleElement(operation, 'type'), tupleElement(operation, 'category'), tupleElement(operation, 'group')), (NULL, 0, 0, 0)) AS operation,
    level,
    category,
    event,
    if(tupleElement(error, 'code') > 0, error, (0, '', '', 0, '', (0, '', '', 0, ''))) AS error,
    message,
    fields,
    JSONExtract(ifNull(fields, ''), 'Map(LowCardinality(String), String)') AS field_map
FROM startup_app_manager_logdb.log_queue;
