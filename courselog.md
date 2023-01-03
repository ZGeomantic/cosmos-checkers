## Step 1: [create stored game]

### 1. 创建一个message
```
$ ignite scaffold message createPost title body
```

### 2. 创建 rules 文件

```
mkdir x/checkers/rules
curl https://raw.githubusercontent.com/batkinson/checkers-go/a09daeb1548dd4cc0145d87c8da3ed2ea33a62e3/checkers/checkers.go | sed 's/package checkers/package rules/' > x/checkers/rules/checkers.go
```


### 3. 创建全局 single Store 对象 systemInfo
```
 ignite scaffold single systemInfo nextId:uint \
    --module checkers \
    --no-message
```

### 4. 创建 map Store 对象

```
ignite scaffold map storedGame board turn black red \
    --index index \
    --module checkers \
    --no-message
```

### 5. GenesisState.SystemInfo 设为 not null
```
// 为第二个字段加上标签 [(gogoproto.nullable) = false]
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  SystemInfo systemInfo = 2 [(gogoproto.nullable) = false];  


$ ignite generate proto-go
```
### 6. 为 StoredGame 添加帮助函数

---

## Step 2: [create message]
### 7. 创建一个 message 对象
```
ignite scaffold message createGame black red \
    --module checkers \
    --response gameIndex
```


### 8. 在 keeper_test 下添加简单的测试用例

---

## Step 3: [implement a msg handler]

### 9. 实现 msgServer.CreateGame 逻辑

### 10. 为 msgServer.CreateGame 补充测试代码

---

## Step 4: [make a move]


### 11. 创建一个新的 message playMove

```
ignite scaffold message playMove gameIndex fromX:uint fromY:uint toX:uint toY:uint \
    --module checkers \
    --response capturedX:int,capturedY:int,winner
```

### 12. 实现 msgServer.PlayMove  以及测试代码

**Note** 并没有把教程中完整的测试用例都粘贴过来

---

## Step 5: [emit an event]
### 13. 代码中通过 EventManager 发出 event，以及测试代码覆盖

主要是测试用例的编写，掌握对```context```中的 event 进行获取的技巧
- 对于 event 的捕捉，有固定的方式``` sdk.StringifyEvents(ctx.EventManager().ABCIEvents())```
- ```StringEvent``` 是先根据 Type 进行分组，再按照执行顺序对数组 ```StringEvent.Attributes```  进行追加的

---


## Step 6: [reject a game]

### 14. 定义一个 rejectGame 对象以及pb
```
ignite scaffold message rejectGame gameIndex --module checkers

# 修改 pb 后重新编译一下
# message StoredGame {
#     ...
#     uint64 moveCount = 6;
# } 
ignite generate proto-go
```


### 15. 实现 msgServer.RejectGame 以及测试代码

---

## Step 7: [put game in order]
### 16. 为游戏实现 FIFO 排序

准备环境变量
```
export alice=$(checkersd keys show alice -a)
export bob=$(checkersd keys show bob -a)

```

启动链
```
ignite chain serve --reset-once
```

执行测试命令
```
checkersd query checkers show-system-info

checkersd tx checkers create-game $alice $bob --from $bob
checkersd query checkers show-system-info
checkersd query checkers show-stored-game 1


checkersd tx checkers create-game $alice $bob --from $bob
checkersd query checkers show-system-info

checkersd tx checkers play-move 2 1 2 2 3 --from $alice
checkersd query checkers show-system-info

```


---
## Step 8: [keep an up-to-date game deadline]

### 17. 新增 deadline 字段，补充相应逻辑

---

## Step 9: [record a winner]

### 18. 添加 winner (test case broken)

还需要给测试用例补上 ```Winner:    "*"```，留在最后一次性修复测试用例吧

---


## Step 10: [auto expiring games]

### 19. 使用 EndBlock 机制实现过期游戏的清理

1.  添加新的逻辑。```AppModule.EndBlock()```函数里添加，属于具体的模块 ```x/checkers/module.go```
   
2.  在 ```app/app.go``` 里注册全局模块的调用顺序
```
app.mm.SetOrderEndBlockers(
        crisistypes.ModuleName,
        ...
+      checkersmoduletypes.ModuleName,
    )
```    


## Step 11: [set a wager]

### 20. 增加 wager 字段，以及 cobra 命令行工具调整

### 21. 为 checker 模块添加可以操作账户资金的权限（bank capability）

### 22. 实现下注、发放奖金、退还资金的操作

### 23. 使用 gomock 完善 bank keeper 的单元测试


### 24. 改造代码，准备可用于集成测试的 testSuite

这里从 cosmos-sdk 和 simapp 中移植了不少的代码

### 25. 利用上一步的改造成果，添加集成测试用例


---
## Step 12: [incentivize players]

### 26. 通过 ```ctx.GasMeter()``` 进行 Gas 费用的收取和退回操作


--- 
## Step 13: [help find a correct move]

### 27. 构造一个查询接口

```
ignite scaffold query canPlayMove gameIndex player fromX:uint fromY:uint toX:uint toY:uint \
    --module checkers \
    --response possible:bool,reason
```

[create stored game]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/3-stored-game.html#some-initial-thoughts
[create message]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/4-create-message.html
[implement a msg handler]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/5-create-handling.html
[make a move]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/6-play-game.html
[emit an event]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/7-events.html
[reject a game]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/8-reject-game.html
[put game in order]: https://interchainacademy.cosmos.network/hands-on-exercise/2-ignite-cli-adv/1-game-fifo.html#
[keep an up-to-date game deadline]: https://interchainacademy.cosmos.network/hands-on-exercise/2-ignite-cli-adv/2-game-deadline.html
[record a winner]: https://interchainacademy.cosmos.network/hands-on-exercise/2-ignite-cli-adv/3-game-winner.html
[auto expiring games]: https://interchainacademy.cosmos.network/hands-on-exercise/2-ignite-cli-adv/4-game-forfeit.html
[set a wager]: https://interchainacademy.cosmos.network/hands-on-exercise/2-ignite-cli-adv/5-game-wager.html
[incentivize players]: https://interchainacademy.cosmos.network/hands-on-exercise/2-ignite-cli-adv/6-gas-meter.html
[help find a correct move]: https://interchainacademy.cosmos.network/hands-on-exercise/2-ignite-cli-adv/7-can-play.html