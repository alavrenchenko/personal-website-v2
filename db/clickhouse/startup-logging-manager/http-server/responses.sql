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

-- ../db/clickhouse/common/http-server/responses.sql

CREATE TABLE IF NOT EXISTS startup_logging_manager_http_server.responses

CREATE TABLE IF NOT EXISTS startup_logging_manager_http_server.response_queue
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'startup_logging_manager.http_server.responses',
    kafka_group_name = 'startup_logging_manager_http_server_responses_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/net/http/server/response_info:ResponseInfo';

CREATE MATERIALIZED VIEW IF NOT EXISTS startup_logging_manager_http_server.response_consumer TO startup_logging_manager_http_server.responses AS
SELECT
    ...
FROM startup_logging_manager_http_server.response_queue;

-------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS startup_logging_manager_http_server.responses
(
    id UUID,
    app Tuple(id UInt64, group_id UInt64, version LowCardinality(String), env LowCardinality(String)),
    app_session_id UInt64,
    http_server_id UInt16,
    request_id UUID,
    timestamp DateTime64(6, 'UTC'),
    status_code Int64,
    body_size Int64,
    content_type LowCardinality(String),

    INDEX id_idx id TYPE bloom_filter GRANULARITY 1,
    INDEX app_id_idx tupleElement(app, 'id') TYPE set(0) GRANULARITY 1,
    INDEX app_group_id_idx tupleElement(app, 'group_id') TYPE set(0) GRANULARITY 1,
    INDEX app_session_id_idx app_session_id TYPE set(0) GRANULARITY 1,
    INDEX http_server_id_idx http_server_id TYPE set(0) GRANULARITY 1,
    INDEX status_code_idx status_code TYPE set(0) GRANULARITY 1,
    INDEX body_size_idx body_size TYPE minmax GRANULARITY 1,
    INDEX content_type_idx lowerUTF8(content_type) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, id)
SETTINGS index_granularity=8192;

CREATE TABLE IF NOT EXISTS startup_logging_manager_http_server.response_queue
(
    id String,
    app Tuple(id UInt64, group_id UInt64, version String, env String),
    app_session_id UInt64,
    http_server_id UInt16,
    request_id String,
    timestamp Int64, -- in microseconds
    status_code Int64,
    body_size Int64,
    content_type String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'startup_logging_manager.http_server.responses',
    kafka_group_name = 'startup_logging_manager_http_server_responses_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/net/http/server/response_info:ResponseInfo';

CREATE MATERIALIZED VIEW IF NOT EXISTS startup_logging_manager_http_server.response_consumer TO startup_logging_manager_http_server.responses AS
SELECT
    toUUID(id) AS id,
    app,
    app_session_id,
    http_server_id,
    toUUID(request_id) AS request_id,
    fromUnixTimestamp64Micro(timestamp, 'UTC') AS timestamp,
    status_code,
    body_size,
    content_type
FROM startup_logging_manager_http_server.response_queue;
