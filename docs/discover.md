# discover机制

## discover页面

- WEB页面


![image-20200805114555955](assets/image-20200805114555955.png)

- Templates代码

  位置在resource/templates/discover.tmpl，其核心是一个表单，对应两个text字段，以及加载一个js代码

  ```javascript
  <script src="{{.prefix}}/js/discover.js"></script>
  ```

  其主要是调用/api/discover/hostname/port接口，然后进行页面渲染(在discover.js:discover(hostname，port)函数中实现)

![image-20200805115158378](assets/image-20200805115158378.png)

这里消息的通知是调用resources/public/js/orchestrator.js中的addAlert函数实现，其是在resources/templates/layout.tmpl中的div

```html
<div id="alerts_container"></div>
```

- resource以及layout加载方式

  在go/app/http.go:standardHttp中实现

  ```go
  	m.Use(render.Renderer(render.Options{
  		Directory:       "resources",
  		Layout:          "templates/layout",
  		HTMLContentType: "text/html",
  	}))
  	m.Use(martini.Static("resources/public", martini.StaticOptions{Prefix: config.Config.URLPrefix}))
  ```
指定加载templates文件、js文件

## discover API

### URL和Handler的映射关系

在go/http/api.go中，支持同步的discover以及异步的discover

```go
func (this *HttpAPI) RegisterRequests(m *martini.ClassicMartini) {
	...
	this.registerAPIRequest(m, "discover/:host/:port", this.Discover)
	this.registerAPIRequest(m, "async-discover/:host/:port", this.AsyncDiscover)
	...
	}
```



异步的discovery就是开启一个协程调用Discover，所以主要还是研究Discover方法。



### Discover方法

#### 整体过程

```go
// Discover issues a synchronous read on an instance
func (this *HttpAPI) Discover(params martini.Params, r render.Render, req *http.Request, user auth.User) {
	//权限判断，以及相关异常处理
	isAuthorizedForAction(req, user)

	//将参数host以及port构建成一个实例，里面涉及域名解析判断
	instanceKey, err := this.getInstanceKey(params["host"], params["port"])

  //根据instancekey采集MySQL相关状态信息，并将相关结果异步写入到orchestrator的后端
	instance, err := inst.ReadTopologyInstance(&instanceKey)

  //如果配置raft，将discover命令使用raft管理起来
	if orcraft.IsRaftEnabled() {
		orcraft.PublishCommand("discover", instanceKey)
	} else {
	  //构建主从关系
		logic.DiscoverInstance(instanceKey)
	}
	//返回成功
	Respond(r, &APIResponse{Code: OK, Message: fmt.Sprintf("Instance discovered: %+v", instance.Key), Details: instance})
}
```

#### ReadTopologyInstance

执行如下SQL并把相关结果写入到后端数据库中

```sql
   show global status like 'Uptime'
   select @@global.hostname, ifnull(@@global.report_host, ''), @@global.server_id, @@global.version, @@global.version_comment, @@global.read_only, @@global.binlog_format, @@global.log_bin, @@global.log_slave_updates
   show slave status
   show global variables like 'rpl_semi_sync_%'
   show master status
   show global status like 'rpl_semi_sync_%'
   select @@global.gtid_mode, @@global.server_uuid, @@global.gtid_executed, @@global.gtid_purged, @@global.master_info_repository = 'TABLE', @@global.binlog_row_image
   show slave hosts
   select count(*) > 0 and MAX(User_name) != '' from mysql.slave_master_info
   select substring_index(host, ':', 1) as slave_hostname from information_schema.processlist where command IN ('Binlog Dump', 'Binlog Dump GTID')
   SELECT SUBSTRING_INDEX(@@hostname, '.', 1)
```

#### Raft机制

- 发布命令

  如果开启raft配置，这里只是记录一个command，PublishCommand将参数Json序列后，记录到raft中

  ```go
  // PublishCommand will distribute a command across the group
  func PublishCommand(op string, value interface{}) (response interface{}, err error) {
  	if !IsRaftEnabled() {
  		return nil, RaftNotRunning
  	}
  	b, err := json.Marshal(value)
  	if err != nil {
  		return nil, err
  	}
  	return store.genericCommand(op, b)
  }
  ```

- 回放-Applier

raft中有回放命令的机制，执行这里操作，对应代码位置是go/logic/command_applier.go

```go
func (applier *CommandApplier) ApplyCommand(op string, value []byte) interface{} {
	switch op {
	...
	case "discover":
		return applier.discover(value)
```

这里的discover对应的是将json反序列，然后调用DiscoverInstance。

```go
func (applier *CommandApplier) discover(value []byte) interface{} {
	instanceKey := inst.InstanceKey{}
	if err := json.Unmarshal(value, &instanceKey); err != nil {
		return log.Errore(err)
	}
	DiscoverInstance(instanceKey)
	return nil
}
```



#### DiscoverInstance

构建主从关系记录到后端数据库中
