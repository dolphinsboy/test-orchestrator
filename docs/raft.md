### 基本点

orchestrator使用的是github.com/hashicorp/raft ，对应的主体代码如下：

```
cd github.com/openark/orchestrator/go/raft
tree -L 2
.
├── applier.go
├── file_snapshot.go
├── fsm.go
├── fsm_snapshot.go
├── http_client.go
├── raft.go
├── rel_store.go
├── snapshot.go
└── store.go

```

### 单机版raft配置



```json
  "RaftEnabled": true,
  "RaftDataDir": "/usr/local/orchestrator/data",
  "RaftBind": "192.168.11.2xx",
  "DefaultRaftPort": 10008,
  "RaftNodes": [
    "192.168.11.2xx"
  ]
```



### 使用单机版raft

raft log对应存在在一个sqlite3的文件中

```sql
[guosong@tdxy-paas-test6 data]$ sqlite3 raft_store.db
SQLite version 3.7.17 2013-05-20 00:56:22
Enter ".help" for instructions
Enter SQL statements terminated with a ";"
sqlite> .tables
raft_log    raft_store
sqlite> select * from raft_log limit 2;
1|1|2|��192.168.11.2xx:10008
2|1|0|{"op":"leader-uri","value":"Imh0dHA6Ly8xOTIuMTY4LjExLjIwMjo4MDAwIg=="}
sqlite> select * from raft_store;
1|CurrentTerm|
2|LastVoteTerm|
3|LastVoteCand|192.168.11.2xx:10008
```

### 启动逻辑

对应的表结构信息：

```sql
sqlite> .schema raft_store
CREATE TABLE raft_store (
			store_id integer,
			store_key varbinary(512) not null,
			store_value blob not null,
			PRIMARY KEY (store_id)
		);
CREATE INDEX store_key_idx_raft_store ON raft_store (store_key)
	;
sqlite> .schema raft_log
CREATE TABLE raft_log (
			log_index integer,
			term bigint not null,
			log_type int not null,
			data blob not null,
			PRIMARY KEY (log_index)
		);
sqlite>
```



在go/app/http.go:standardHttp

在raft.go:Setup中会创建raft。




### 文档参考

- [基于hashicorp/raft的分布式一致性实战教学](<https://zhuanlan.zhihu.com/p/58048906>)


