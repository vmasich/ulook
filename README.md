# URL lookup service

## Problem statement

We have an HTTP proxy that is scanning traffic looking for malware
URLs. Before allowing HTTP connections to be made, this proxy asks a
service that maintains several databases of malware URLs if the
resource being requested is known to contain malware.

Write a small web service, in the language/framework your choice, that
responds to GET requests where the caller passes in a URL and the
service responds with some information about that URL. The GET
requests look like this:

```sh
GET /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}
```

The caller wants to know if it is safe to access that URL or not. As
the implementer you get to choose the response format and
structure. These lookups are blocking users from accessing the URL
until the caller receives a response from your service.

The service must be containerized.

Give some thought to the following:

1. The size of the URL list could grow infinitely, how might you
scale this beyond the memory capacity of this VM? Bonus if you
implement this.

2. The number of requests may exceed the capacity of this VM, how might
you solve that? Bonus if you implement this.

3. What are some strategies you might use to update the service with new
URLs? Updates may be as much as 5 thousand URLs a day with updates
arriving every 10 minutes.

## Usage

The example storage cluster consists of 2 nodes.
More can be added - need to create volumes and assign sharding filter ranges.

```sh
docker network create mynet
docker volume create store1
docker volume create store2

docker run --rm  -d --net mynet --name nats  nats

docker run --rm -d  --net mynet -v store1:/db --name dbstore1 urllookup /dbstore -filter="0m"

docker run --rm -d  --net mynet -v store2:/db --name dbstore2 urllookup /dbstore -filter="m~"

docker run --rm -d -p 3333:3333 --net mynet --name restapi urllookup /restapi
```


## Architecture

### The implementation uses microservices architecture.

Lookup services uses following packages:

- BoltDb https://github.com/boltdb/bolt. It is being used as storage
  engine in ETCD, InfluxDb, Prometheus

- NATS messaging system https://nats.io

- Gin HTTP framework (relatively new framework with Go Context package support)

### The components of a the URL Lookup service are

- HTTP REST fronted

- NATS messaging system

- BoltStore - distributed URL data store.


### Alternatives and the rationale:

Alternatively, some existing database can be used as a backend.

The good candidates can be DynamoDB, Dgraph, CockroachDB, Cassandra,
etc.

However implementing something from scratch is a fun - and it is not
that complex.

### The intended setup:

**HTTP REST frontend** should be run as a multiple instances, depending
 on requirements behind Load Balancer ( Kubernetes, NGNX, AWS Elastic
 LB, etc. )

**NATS messaging** system should be run as a clustered setup on the low
 latency network

**Boltstore** instances should be run as a sharded setup with multiple
nodes. Each instance serves a range of URLs , defined as a range (in
current very simple implementation), i.e. for 2 store nodes [0k] for
the first, [k~]* for the second store node. In production
implementation of the filters can be adaptive, the data on the nodes
can be replicated, migrated, split, or merged.  The current max size
of the node is 1TB.

- The implemented architecture is easily horizontally scalable and addresses
  pp. 1 and 2 in the problem statement.

- It is designed with a potential for the future
  optimization. Additional *Data Store types* can be added, for
  example a *Cache with Hit counters* for the most repeated lookups,
  and "Bloom Filters" which will be even faster than Cache.

- The *Cache Nodes* can be prepared as a background process using data
  from the Data Store, or on the specialized nodes for the counters

- No architectural changes are required for that, because the response
  from the fastest node will be returned to the user.

- Another optimization should be implemented in the Data
  Store. Current internal structure is flat, so buckets are in the
  same level, the optimized implementation will use hierarchical
  structure with nested buckets, following the Domain Names levels and
  path components. It will speed up the lookups. It is not implemented
  because current implementation complexity is already beyond the scope of this
  "staightforward" exersize.

- To address pp. 3 in the problem statement, the REST API supports
  bulk URL updates. JSON is used for the simplicity, CSV will be more efficient.
	However 5000 URL a day  sounds too little -  may be 5000 per request?


## API Reference

### URL Lookup
```sh
GET /urlinfo/1/:hostname/:path
```
#### Parameters

hostname: hostname and port

path: original path and query string

#### Payload

No payload

#### Response

HTTP status codes

200: URL found

404: URL not found

504: Timeout

- There is no Response Body to reduce amount of data transferred

### Bulk URL update
```sh
POST /urlinfo/bulkupdate
```
#### Parameters

None

#### Headers

```sh
Content-Type: application/json
```
#### Payload

Used JSON for the siplicity of implementation. CSV will be more efficient

- example
```json
[
    {
        "op": "+",
        "h": "abc.u",
        "pq": "tadam?r=2"
    },
    {
        "op": "+",
        "h": "pvc.b.c",
        "pq": "tadam"
    },
]
```
