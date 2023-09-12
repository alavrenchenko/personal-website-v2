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

-- ../db/clickhouse/actiondb/transactions.sql

CREATE TABLE IF NOT EXISTS test_actiondb.transactions
(
    id UUID,
    app Tuple(id UInt64, group_id UInt64, version LowCardinality(String), env LowCardinality(String)),
    app_session_id UInt64,
    created_at DateTime64(6, 'UTC'),
    start_time DateTime64(6, 'UTC'),

    INDEX id_idx id TYPE bloom_filter GRANULARITY 1,
    INDEX app_id_idx tupleElement(app, 'id') TYPE set(0) GRANULARITY 1,
    INDEX app_group_id_idx tupleElement(app, 'group_id') TYPE set(0) GRANULARITY 1,
    INDEX app_session_id_idx app_session_id TYPE set(0) GRANULARITY 1
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(created_at)
ORDER BY (created_at, id)
SETTINGS index_granularity=8192;

CREATE TABLE IF NOT EXISTS test_actiondb.transaction_queue
(
    id String,
    app Tuple(id UInt64, group_id UInt64, version String, env String),
    app_session_id UInt64,
    created_at Int64, -- in microseconds
    start_time Int64 -- in microseconds
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'localhost:9092',
    kafka_topic_list = 'testing.base.transactions',
    kafka_group_name = 'test_base_transactions_clickhouse',
    kafka_format = 'ProtobufSingle',
    kafka_schema = 'personalwebsite/actions/transaction:Transaction';

CREATE MATERIALIZED VIEW IF NOT EXISTS test_actiondb.transaction_consumer TO test_actiondb.transactions AS
SELECT
    toUUID(id) AS id,
    app,
    app_session_id,
    fromUnixTimestamp64Micro(created_at, 'UTC') AS created_at,
    fromUnixTimestamp64Micro(start_time, 'UTC') AS start_time
FROM test_actiondb.transaction_queue;
