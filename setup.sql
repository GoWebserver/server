create role server
            with login = true;

create keyspace server with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
use server;

create table access
(
    id             timeuuid primary key,
    code           smallint,
    duration       int,
    error          text,
    https          boolean,
    method         text,
    searchduration int,
    uri            text,
    writeerr       text
)
    with caching = {'keys': 'ALL', 'rows_per_partition': 'ALL'}
     and compaction = {'class': 'SizeTieredCompactionStrategy'}
     and compression = {'sstable_compression': 'org.apache.cassandra.io.compress.LZ4Compressor'}
     and dclocal_read_repair_chance = 0
     and speculative_retry = '99.0PERCENTILE';

create table apiaccess
(
    id       timeuuid primary key,
    duration int,
    error    text,
    request  text
)
    with caching = {'keys': 'ALL', 'rows_per_partition': 'ALL'}
     and compaction = {'class': 'SizeTieredCompactionStrategy'}
     and compression = {'sstable_compression': 'org.apache.cassandra.io.compress.LZ4Compressor'}
     and dclocal_read_repair_chance = 0
     and speculative_retry = '99.0PERCENTILE';

create table mime
(
    extension text primary key,
    mimetype  text
)
    with caching = {'keys': 'ALL', 'rows_per_partition': 'ALL'}
     and compaction = {'class': 'SizeTieredCompactionStrategy'}
     and compression = {'sstable_compression': 'org.apache.cassandra.io.compress.LZ4Compressor'}
     and dclocal_read_repair_chance = 0
     and speculative_retry = '99.0PERCENTILE';

