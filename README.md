## test-orchestrator
test-orchestrator with go mod

使用go mod模式管理包，原来orchestrator使用govendor进行管理

## http

### 框架
使用[martini框架](https://github.com/go-martini/martini)进行开发，不过现在martini框架已经不再维护。

之前看过的open-falcon项目使用的是[gin-gonic](https://github.com/gin-gonic/gin)，自己开发项目的时候也可以选择。

### orchestrator

支持API以及WEB两种模式，WEB调用API接口，开发比较方面。open-falcon项目使用gin-gonic开发API接口，WEB界面使用Python Flask开发。


## Discover机制

具体参照[docs/discover](https://github.com/dolphinsboy/test-orchestrator/blob/master/docs/discover.md)

## Raft

具体参照说明[docs/raft](https://github.com/dolphinsboy/test-orchestrator/blob/master/docs/raft.md)