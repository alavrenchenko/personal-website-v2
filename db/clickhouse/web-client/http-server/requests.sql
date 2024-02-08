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

-- ../db/clickhouse/common/http-server/requests.sql

CREATE TABLE IF NOT EXISTS web_client_http_server.requests

CREATE TABLE IF NOT EXISTS web_client_http_server.request_queue
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'web_client.http_server.requests',
    kafka_group_name = 'web_client_http_server_requests_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/net/http/server/request_info:RequestInfo';

CREATE MATERIALIZED VIEW IF NOT EXISTS web_client_http_server.request_consumer TO web_client_http_server.requests AS
SELECT
    ...
FROM web_client_http_server.request_queue;

-------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS web_client_http_server.requests
(
    id UUID,
    app Tuple(id UInt64, group_id UInt64, version LowCardinality(String), env LowCardinality(String)),
    app_session_id UInt64,
    http_server_id UInt16,
    status Enum8('unspecified' = 0, 'new' = 1, 'in_progress' = 2, 'success' = 3, 'failure' = 4),
    start_time DateTime64(6, 'UTC'),
    end_time Nullable(DateTime64(6, 'UTC')),
    elapsed_time_us Nullable(Int64), -- in microseconds
    url String,
    method LowCardinality(String),
    protocol LowCardinality(String),
    host LowCardinality(String),
    remote_addr String, -- "IP:port"
    request_uri String,
    content_length Int64,
    headers Map(LowCardinality(String), Array(String)),
    x_real_ip String,
    x_forwarded_for String,
    content_type LowCardinality(String),
    origin String,
    referer String,
    user_agent String,
    _version_stamp UInt8,

    INDEX id_idx id TYPE bloom_filter GRANULARITY 1,
    INDEX app_id_idx tupleElement(app, 'id') TYPE set(0) GRANULARITY 1,
    INDEX app_group_id_idx tupleElement(app, 'group_id') TYPE set(0) GRANULARITY 1,
    INDEX app_session_id_idx app_session_id TYPE set(0) GRANULARITY 1,
    INDEX http_server_id_idx http_server_id TYPE set(0) GRANULARITY 1,
    INDEX status_idx status TYPE set(0) GRANULARITY 1,
    INDEX elapsed_time_us_idx elapsed_time_us TYPE minmax GRANULARITY 1,
    INDEX url_idx lowerUTF8(url) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX method_idx method TYPE set(0) GRANULARITY 1,
    INDEX remote_addr_idx remote_addr TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX content_length_idx content_length TYPE minmax GRANULARITY 1,
    INDEX headers_key_idx arrayMap(k -> lowerUTF8(k), mapKeys(headers)) TYPE set(0) GRANULARITY 1,
    INDEX headers_value_idx arrayMap(v -> lowerUTF8(v), flatten(mapValues(headers))) TYPE bloom_filter GRANULARITY 1,
    INDEX x_real_ip_idx x_real_ip TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX x_forwarded_for_idx x_forwarded_for TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX content_type_idx lowerUTF8(content_type) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX origin_idx lowerUTF8(origin) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX referer_idx lowerUTF8(referer) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX user_agent_idx lowerUTF8(user_agent) TYPE tokenbf_v1(8192, 3, 0) GRANULARITY 1,
    INDEX _version_stamp_idx _version_stamp TYPE set(0) GRANULARITY 1
)
ENGINE = ReplacingMergeTree(_version_stamp)
PARTITION BY toYYYYMM(start_time)
ORDER BY (start_time, id)
SETTINGS index_granularity=8192;

CREATE TABLE IF NOT EXISTS web_client_http_server.request_queue
(
    id String,
    app Tuple(id UInt64, group_id UInt64, version String, env String),
    app_session_id UInt64,
    http_server_id UInt16,
    status Enum8('unspecified' = 0, 'new' = 1, 'in_progress' = 2, 'success' = 3, 'failure' = 4),
    start_time Int64, -- in microseconds
    end_time Nullable(Int64), -- in microseconds
    elapsed_time_us Nullable(Int64), -- in microseconds
    url String,
    method String,
    protocol String,
    host String,
    remote_addr String,
    request_uri String,
    content_length Int64,
    headers String,
    x_real_ip String,
    x_forwarded_for String,
    content_type String,
    origin String,
    referer String,
    user_agent String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'web_client.http_server.requests',
    kafka_group_name = 'web_client_http_server_requests_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/net/http/server/request_info:RequestInfo';

CREATE MATERIALIZED VIEW IF NOT EXISTS web_client_http_server.request_consumer TO web_client_http_server.requests AS
SELECT
    toUUID(id) AS id,
    app,
    app_session_id,
    http_server_id,
    status,
    fromUnixTimestamp64Micro(start_time, 'UTC') AS start_time,
    fromUnixTimestamp64Micro(end_time, 'UTC') AS end_time,
    elapsed_time_us,
    url,
    method,
    protocol,
    host,
    remote_addr,
    request_uri,
    content_length Int64,
    JSONExtract(headers, 'Map(LowCardinality(String), Array(String))') AS headers,
    x_real_ip,
    x_forwarded_for,
    content_type,
    origin,
    referer,
    user_agent,
    if(end_time IS NOT NULL, 2, 1) AS _version_stamp
FROM web_client_http_server.request_queue;
