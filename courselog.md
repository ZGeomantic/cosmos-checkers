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

## Step3: [implement a msg handler]

### 9. 实现 msgServer.CreateGame 逻辑

### 10. 为 msgServer.CreateGame 补充测试代码

---

## Step4: [make a move]


### 11. 创建一个新的 message playMove

```
ignite scaffold message playMove gameIndex fromX:uint fromY:uint toX:uint toY:uint \
    --module checkers \
    --response capturedX:int,capturedY:int,winner
```

### 12. 
---

[create stored game]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/3-stored-game.html#some-initial-thoughts
[create message]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/4-create-message.html
[implement a msg handler]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/5-create-handling.html

[make a move]: https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/6-play-game.html