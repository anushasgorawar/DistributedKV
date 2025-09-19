Final expectation:

`(base) anushasg@Anushas-MacBook-Air bin % curl 'http:127.0.0.5:8080/get?key=party'`      

Shard = 3, current shard = 3, addr = "127.0.0.5:8080", Value = "Great", error = <nil>%

`(base) anushasg@Anushas-MacBook-Air bin % curl 'http:127.0.0.55:8080/get?key=party'`

Shard = 3, current shard = 3, addr = "127.0.0.5:8080", Value = "Great", error = <nil>%

`(base) anushasg@Anushas-MacBook-Air bin % curl 'http:127.0.0.55:8080/get?key=health'`

Shard = 3, current shard = 3, addr = "127.0.0.5:8080", Value = "", error = <nil>%     

`(base) anushasg@Anushas-MacBook-Air bin % curl 'http:127.0.0.55:8080/set?key=health&value=weak'`

Error = read-only mode, shardIdx = 3, current shard = 3%    

`(base) anushasg@Anushas-MacBook-Air bin % curl 'http:127.0.0.5:8080/set?key=health&value=weak' `

Error = <nil>, shardIdx = 3, current shard = 3%

`(base) anushasg@Anushas-MacBook-Air bin % curl 'http:127.0.0.5:8080/get?key=health' `          

Shard = 3, current shard = 3, addr = "127.0.0.5:8080", Value = "weak", error = <nil>%  `anushasg@Anushas-MacBook-Air bin % curl 'http:127.0.0.55:8080/get?key=health'`

Shard = 3, current shard = 3, addr = "127.0.0.5:8080", Value = "weak", error = <nil>% 



## Set up

In zsh terminal:

1. `go mod init github.com/anushasgorawar/DistributedKV` to initialise



2. link to boltDB: https:github.com/boltdb/bolt

`go get github.com/boltdb/bolt/...` to download boltdb source code into the module

bolt command line utility into your $GOBIN path.



Why BoltDB? It's an embedded key value db - it runs as a process instead of being a seperate server. Data is Key values are stored in Buckets.

It writes to disk-SSD



3. once `db, err := bolt.Open("my.db", 0600, nil)` is added and run,

we get my.db file created.



4. define flags(whta be pass when running the program like ./file.text --run=dry) for db like location. Package flag implements command-line flag parsing.

Need to parse the flags using `flag.Parse()` before we can use the flags

`go run main.go --db-location="my.db --http-address="127.0.0.1:8080"`



5. define what people can do with our service. With HTTP API, HTTP handlers



6. create a package - db, to isolate server and db

our type - database

After writing that module:

add this in import

`"github.com/anushasgorawar/DistributedKV/db"`

run `go mod tidy`



7. closeFunc = boltDB.Close

We are assigning a method(boltDB.Close) to a closeFunc



8. SetKey func: sets a key value pair

There's no direct function so we've to add our own transactions.



9. cleanUpFunc -> boltdb.close return nil if no error and closed properly. So, if it is not giving an error, it'll not run it again. But, for simplicity, removed clean up.



Till now, I set up a webserver with two endpoints "/get" and "/set" which can retrieve, add and update a key value in the. bolt db.



## Sharding

1. Determine which key is stored in which shard, our sharding is static. Hence we have have one toml file.



2. go library for toml parsing: https:github.com/BurntSushi/toml

`go get github.com/BurntSushi/toml@latest`

added `"github.com/BurntSushi/toml"` to import block

read from sharding.toml using os package OR better, toml.decode()



3. config package has config related code. Each shard has a unique set of keys. 



4. Error: 

```

2025/09/02 00:07:34 Unable to decode config File sharding.toml, error: toml: line 10 (last key "shard"): type mismatch for config.Shard: expected table but found []map[string]any

exit status 1

```

config struct needs to have a map and not just a shard type

To check if sharding is correct: log.Printf("%#v", &c)

2025/09/02 00:11:28 &config.Config{Shard:[]config.Shard{config.Shard{name:"", idx:0}, config.Shard{name:"", idx:0}, config.Shard{name:"", idx:0}, config.Shard{name:"", idx:0}}}

fix: Capitalised the first letter of keys Name and Idx

err: 

go run . --db-location="my.db" --config-file="sharding.toml"

2025/09/06 20:32:53 &config.Config{Shard:[]config.Shard(nil)}



5.  In toml [[ xyz ]] describes an array object in an array

The TOML library is case-insensitive and uses some loose matching rules:

"name" in TOML maps to Name in Go

"idx" in TOML maps to Idx in Go



6. We pass the shard also as a flag. We check for that shard's index using the shard name. 



7. Each shard will have it's sharding key which is unique.

hash(key) % total shards = current index



8. Hash function - package fnv.

Use this in /set. 

hash.Sum64() will give a 64 bit hash number.



9. addr in sharding.toml is whats accessible from every other shard.

addr in script is the socket we are listening to



10. get url mismatch

> curl 'http:127.0.0.1:8082/get?key=boy' 

value="", error:<nil>%

>curl 'http:127.0.0.1:8081/get?key=boy' 

value="praj", error:<nil>%  



11. 

> curl 'http:127.0.0.1:8081/set?key=chair&value=pink'

shardInd=0, value="pink" error:<nil>%

> curl 'http:127.0.0.1:8083/set?key=chair&value=pink'

shardInd=0, value="pink", error:<nil>% 

> curl 'http:127.0.0.1:8081/get?key=chair'           

value="pink", error:<nil>% 

> curl 'http:127.0.0.1:8082/get?key=chair'

value="", error:<nil>%



12. Basically no matter which addr we call[9] in our curl command, the key value should be set to and get from the hash function. 



13. Move the sharding related code to config.

handle calls get function with the correct url



14. If a shard name "potato" is not there, we get a log: Shard "potato" not found.



## testing



1. To run all tests:

> go test ./...       

?       github.com/anushasgorawar/DistributedKV [no test files]

?       github.com/anushasgorawar/DistributedKV/db      [no test files]

?       github.com/anushasgorawar/DistributedKV/web     [no test files]

ok      github.com/anushasgorawar/DistributedKV/config  (cached)



2. TestWebServer: We dont know the address of the server so we add 2 variables instead.

create map for address

compute shard index -simulation before the actual code

log.Printf("adjective: %d", s.GetShard("adjective")) 0

log.Printf("adverb: %d", s.GetShard("adverb")) 1

t.Fail()

adjectiveShard := 0

adverbShard = 1

keys := map[string]int{

	"adjective": 0,

	"adverb":    1,

}



3. Go through testing of webserver



## Sharding using powers of 2



1. resharding in powers of 2 - While rebalancing shards, the keys will not move across existing servers.



2. If we delete the shard without resharding, the data is lost. 

I deleted the animal and thing shard from sharding toml file and when I ran the curl command that had a key in "thing" shard was not available as the new shard of that key was 2.



3. populate.sh has a script to populate the shards[small and managable pieces of a db].

If we are copying dbs, we need to stop the application.



4. We copy dbs 

```

cp name.db animal.db 

cp place.db thing.db 

```

5. Once copies, we can delete the keys that don't belong in these shards.

For example, A key called "watermelon" was present in shard 1 and when the shards increased to 4 from 2, we copied the keys of shard 1 to shard 3 and 0 to 2.

If watermelon key now belongs to shard 3, we can remove it from shard 1.



6. Code to delete keys that don't belong to a shard.

We can't modify a collection while we iterate.

Hence 2 loops



7. PurgeHandler deletes all the keys that dont belong in that shard.



8. if d.shards.CurrInd != d.shards.GetShard, delete that key.



test: We have 2 shards. key-xyz is in shard 1. We have 4 shards now. key-xyz's new hash assigns it shard 3. It should be moved to shard 3 and deleted in shard 1. We can call /purge endpoint to delete extra keys in shard 1 but we need something to re-shard keys in shard 3

We won't be able to fetch the value of the key as it'll be handled by a differnt address.



## Benchmarking writes in key value database



prereq: install fish.



1. When the shard is correct, it takes 1ms but when its different, 14ms i.e when it redirects

[Benchmarking writes means measuring the performance of write operations in a storage system, database, file system, or application.]



2. cmd will be excecuted from time to time. 

we can benchmark using embedded go tools.

We want to becnhmark against our db.



3. Write a random key to our instance in cmd



4. Since instances are running in same server, we can benchmark 1.

when in production, it'll be different servers, we can benchmark all of them.

Since in our application, we redirect to different instance for specifies key, we will eventually benchmark all four of them.



5. Actual benchmarking for 4 writes:

step 1: Avg(t1,t2,t3,t4)

step 2: throughput=total bytes written/total time taken=(1/avg)*(bytes written)



6. can calculate min and max also.

min is around 14ms and max around144ms

avg of 100 writes: 10.495249ms 

avg of 1000 writes: 9.807323ms



Trace and see why it's taking this much time.



Everytime a query, it checks with boltdb. Could be a cause for latency.



We could try go routines and implement concurrency



Concurrency does not improve throughput because boltdb queries are not concurrent.



QPS=Queries Per Second

qps is improved when we write to memeory /dev/shm/name.db instead of local name.db



7. Avergae Time taken for function write: 9.89756ms 

Avergae Time taken for function read: 417.117µs 



8. boltDB.NoSync = true ? Doesnt sync to disc - not to be used in production

- allows to write to boltdb faster - will check performance.

Before: 

func write- avg: 320.94687ms, min: 0s, max: 3.209468416s, QPS: 31.2

After:

func write- avg: 15.34517ms, min: 0s, max: 153.451417ms, QPS: 651.7

Total QPS= 3281.4 

Setting the NoSync flag will cause the database to skip fsync() calls after each commit. This can be useful when bulk loading data into a database and you can restart the bulk load in the event of a system failure or database corruption. Do not set this flag for normal use.



func write- avg: 14.683404ms, min: 0s, max: 146.792959ms, QPS: 681.0

func write- avg: 14.900566ms, min: 0s, max: 149.0055ms, QPS: 671.1

func write- avg: 15.001404ms, min: 0s, max: 150.013917ms, QPS: 666.6

func write- avg: 15.027391ms, min: 0s, max: 150.27375ms, QPS: 665.5

func write- avg: 15.053954ms, min: 0s, max: 150.539375ms, QPS: 664.3

2025/09/16 20:38:24 Total QPS= 3348.4, set 500 keys

func read- avg: 7.318179ms, min: 0s, max: 73.181541ms, QPS: 1366.5

func read- avg: 7.355995ms, min: 0s, max: 73.55975ms, QPS: 1359.4

func read- avg: 7.426425ms, min: 0s, max: 74.22125ms, QPS: 1346.5

func read- avg: 7.431091ms, min: 0s, max: 74.31075ms, QPS: 1345.7

func read- avg: 7.474987ms, min: 0s, max: 74.749709ms, QPS: 1337.8

2025/09/16 20:38:24 Total QPS= 6755.9, set 500 keys



Benchmarking all the shards, four servers

To check single shard: temperorily edit the config to have only one shard.



after adding: io.Copy(io.Discard, resp.Body); it significantly improved the QPS



(base) anushasg@Anushas-MacBook-Air DistributedKV % go run cmd/main.go

func write- avg: 10.673841ms, min: 0s, max: 106.686625ms, QPS: 936.8

func write- avg: 12.390316ms, min: 0s, max: 123.902916ms, QPS: 807.1

func write- avg: 12.583616ms, min: 0s, max: 125.836042ms, QPS: 794.7

func write- avg: 12.5777ms, min: 0s, max: 125.776875ms, QPS: 795.1

func write- avg: 12.637012ms, min: 0s, max: 126.369958ms, QPS: 791.3

2025/09/16 23:57:55 Total QPS= 4124.9, set 500 keys

func read- avg: 4.471741ms, min: 0s, max: 44.717166ms, QPS: 2236.3

func read- avg: 4.520004ms, min: 0s, max: 45.199875ms, QPS: 2212.4

func read- avg: 4.532416ms, min: 0s, max: 45.324083ms, QPS: 2206.3

func read- avg: 4.571033ms, min: 0s, max: 45.710125ms, QPS: 2187.7

func read- avg: 4.585308ms, min: 0s, max: 45.852917ms, QPS: 2180.9

2025/09/16 23:57:55 Total QPS= 11023.5, set 500 keys



io.Discard is a special writer in Go that simply discards anything written to it — it acts like a "black hole" for data.

Why do this?

Reading and discarding the entire response body ensures that the HTTP connection can be reused (in HTTP/1.1 keep-alive connections).

If you don't fully read and close the response body, the connection might not be reused, which can hurt performance.

After this, you should still call resp.Body.Close() to release resources.
Or else, tcp ports will run out.

resp.Body is an io.ReadCloser backed by a network connection.

If you don’t close it, the underlying TCP connection may remain open, leading to resource leaks (open file descriptors, unused sockets).

Over time, this can exhaust system resources and cause your program to fail.



We open connections and close then again. It's not a good practice to close the connections so quickly. lets keep the max connections 32



http is the entire http package. http client is a struct representing http client.

It represents an HTTP client that you can use to make HTTP requests.

It provides more control over HTTP requests than the convenience functions like http.Get(). 

because it lets you customize many aspects of how requests are made and handled. 

customisations like timeouts, Redirect Policy, Custom Headers and Methods[Any method (GET, POST, PUT, etc.)],  Cookie Jar[cookies management], tune connection pooling and concurrency behavior for better performance.



after adding http client instead of http:

func write- avg: 9.228312ms, min: 0s, max: 92.25675ms, QPS: 1083.5

func write- avg: 9.538145ms, min: 0s, max: 95.381291ms, QPS: 1048.4

func write- avg: 9.541704ms, min: 0s, max: 95.416917ms, QPS: 1048.0

func write- avg: 9.621183ms, min: 0s, max: 96.211709ms, QPS: 1039.4

func write- avg: 9.7018ms, min: 0s, max: 97.017875ms, QPS: 1030.7

2025/09/17 00:20:16 Total QPS= 5250.1, set 500 keys

func read- avg: 4.600125ms, min: 0s, max: 46.001041ms, QPS: 2173.8

func read- avg: 4.825495ms, min: 0s, max: 48.254792ms, QPS: 2072.3

func read- avg: 4.9085ms, min: 0s, max: 49.084875ms, QPS: 2037.3

func read- avg: 5.008429ms, min: 0s, max: 50.084166ms, QPS: 1996.6

func read- avg: 5.0101ms, min: 0s, max: 50.100917ms, QPS: 1996.0

2025/09/17 00:20:16 Total QPS= 10276.0, set 500 keys



After changing IPs of the servers: [to correct ones]

[Before I had added one IP different ports, now it's different IPs]

func write- avg: 306.925333ms, min: 0s, max: 3.069131791s, QPS: 32.6

func write- avg: 317.60782ms, min: 0s, max: 3.176077916s, QPS: 31.5

func write- avg: 317.857454ms, min: 0s, max: 3.178574375s, QPS: 31.5

func write- avg: 320.156233ms, min: 0s, max: 3.201562083s, QPS: 31.2

func write- avg: 323.663329ms, min: 0s, max: 3.236633084s, QPS: 30.9

2025/09/17 10:59:54 Total QPS= 157.7, set 500 keys

func read- avg: 5.698579ms, min: 0s, max: 56.985666ms, QPS: 1754.8

func read- avg: 5.93865ms, min: 0s, max: 59.386416ms, QPS: 1683.9

func read- avg: 6.021445ms, min: 0s, max: 60.214375ms, QPS: 1660.7

func read- avg: 6.035991ms, min: 0s, max: 60.359792ms, QPS: 1656.7

func read- avg: 6.033533ms, min: 0s, max: 60.33525ms, QPS: 1657.4

2025/09/17 10:59:54 Total QPS= 8413.6, set 500 keys



QPS for reading goes up with the number of iterations and concurrency.

Write is stable at 155-165 QPS



If bolt db sync is disabled:

(base) anushasg@Anushas-MacBook-Air DistributedKV % go run /Users/anushasg/Documents/projects/DistributedKV/cmd/main.go -iterations=1000 concurrency=10

func write- avg: 78.39702ms, min: 0s, max: 783.916625ms, QPS: 1275.5

func write- avg: 79.449233ms, min: 0s, max: 794.492042ms, QPS: 1258.7

func write- avg: 80.374841ms, min: 0s, max: 803.748041ms, QPS: 1244.2

func write- avg: 80.714637ms, min: 0s, max: 807.146125ms, QPS: 1238.9

func write- avg: 81.239062ms, min: 0s, max: 812.3905ms, QPS: 1230.9

2025/09/17 11:03:11 Total QPS= 6248.2, set 5000 keys

func read- avg: 4.366666ms, min: 0s, max: 43.6665ms, QPS: 2290.1

func read- avg: 4.4089ms, min: 0s, max: 44.088834ms, QPS: 2268.1

func read- avg: 4.419229ms, min: 0s, max: 44.192167ms, QPS: 2262.8

func read- avg: 4.423858ms, min: 0s, max: 44.2385ms, QPS: 2260.5

func read- avg: 4.438487ms, min: 0s, max: 44.38475ms, QPS: 2253.0

2025/09/17 11:03:11 Total QPS= 11334.5, set 5000 keys



## Replication



1. Logical replication and physicall replication and how to manage it.



2. replication: We'll not lose data if a server goes down. Replcation is not backup. backup is for human error. replication stores all the changes from one server to another. If something is deleted in one, it's deleted in another as well.

Replicas can't accept writes, only reads.



3. db name of the replica should be different. otherwise, they'll try to write to the same file.

4. add replica-address to toml file

5. update database struct to have readonly flag. If replica, then readonly. else false.

6. If replica=true: test of set will fail

(base) anushasg@Anushas-MacBook-Air DistributedKV % go test ./db/db_test.go

--- FAIL: TestSetGet (0.03s)

    db_test.go:22: could not set key: "read-only mode"

    db_test.go:31: expected: "Round", recieved: ""

7. defer closeFunc()

	t.Cleanup(func() { closeFunc() })

Use defer when the resource should be closed as soon as the current function exits.

Use t.Cleanup when the resource should be closed at the end of the test (or subtest), even if it was created in a helper.



8. How is replica implemented? First when we set a key, it'll be stored in the main db, not yet replicated. The replica reads this from the amin.db. Once replicated, the keys are deleted in the main db.

Several approaches. Simple one.



9. Write will happen in 2 buckets - main and replica bucket. 

After we are done, we delete values from replica bucket. 



10. A function to tell replica which changes are done since last time - Replicate.

if key and value match, delete it.

if value is different, i.e. new version, return error

DeleteReplicatedKey : deletes key from replica queue if the value matches with the value



Basically, everything in the main db is a queue for replica to be replicated.

If the value is already updated, it's ignored and deleted in the main db.

if not, it's sends an error.





So far:

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/set?key=moon&value=full'

shard=2 addr=127.0.0.4:8080 value="full", error:<nil>%                                                                      

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/get?key=moon'           

shard=2 addr=127.0.0.4:8080 value="full", error:<nil>%                                                                      

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/purge'       

err=<nil>%                                                                                                                  

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/get?key=moon'

shard=2 addr=127.0.0.4:8080 value="full", error:<nil>%                                                                      

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/next-replication-key'

{"Key":"moon","Value":"full","Err":null}

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.44:8080/set?key=moon&value=full'

shard=2 addr=127.0.0.4:8080 value="full", error:read-only mode%  



11. 
(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/set?key=moon&value=full'   
shard=2 addr=127.0.0.4:8080 value="full", error:<nil>%   

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/set?key=mars&value=full'
shard=2 addr=127.0.0.4:8080 value="full", error:<nil>%   

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/next-replication-key'
{"Key":"mars","Value":"full","Err":null}

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/delete-next-replication-key?key=mars&value=full'
ok%                              

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.4:8080/next-replication-key'
{"Key":"moon","Value":"full","Err":null}

12. if its a replica, it starts a replica thread, start a replica client. i.e replica package


13. (c *client) loop()  will contact the master server/leader and downloads the next replictaion key


14. Issue: 
```
(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.22:8080/get?key=mars'
curl: (7) Failed to connect to 127.0.0.22 port 8080 after 0 ms: Couldn't connect to server
(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.33:8080/get?key=mars'
shard=2 addr=127.0.0.4:8080 value="full", error:<nil>%                                                                      
(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.44:8080/get?key=mars'
curl: (7) Failed to connect to 127.0.0.44 port 8080 after 0 ms: Couldn't connect to server
(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.55:8080/get?key=mars'
curl: (7) Failed to connect to 127.0.0.55 port 8080 after 1 ms: Couldn't connect to server
```
Changed sleep from 1 minute to 1 second. The socket was not closing until that time was over.


To replicate:
cp name.db name-replica.db


To test:
1. I have set the key and get it working

(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.55:8080/get?key=earth'
shard=3 addr=127.0.0.5:8080 value="moon", error:<nil>%  

2. I'll delete the 4th db and try the same query.
launched again.
(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.55:8080/get?key=earth'
curl: (7) Failed to connect to 127.0.0.55 port 8080 after 1 ms: Couldn't connect to server
Might have to close that connection before I start

Had to run: killall DistributedKV
(base) anushasg@Anushas-MacBook-Air DistributedKV % curl 'http://127.0.0.55:8080/get?key=earth'
shard=3 addr=127.0.0.5:8080 value="moon", error:<nil>%  
Works!