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

### 27. 构造一个查询接口 canPlayMove

```
ignite scaffold query canPlayMove gameIndex player fromX:uint fromY:uint toX:uint toY:uint \
    --module checkers \
    --response possible:bool,reason
```

### 28. 实现查询接口的逻辑，以及单元、集成测试的补充

```
 ignite chain build
```

---

## Step 14: [play with IBC tokens]

### 29. 调整proto，在创建游戏时指定 denom, 由 StoredGame.GetWagerCoin() 根据入参指定下注的币种


### 30. 对 IBC 进行测试，为测试函数都加上 denom 入参。调通了多个币都可以进行质押的测试用例

这里还是有点奇怪，这个算是 IBC 吗？不是在同一个链上添加多个币种作为资产而已吗？
并且发起IBC转账的时候，其实不用特殊的函数吗？仍然用 bank.SendCoinsFromAccountToModule or BaseSendKeeper.SendCoins 就可以处理 IBC 之间的 token 转移？
It's still a bit strange here, is this an IBC? Isn't it just adding multiple currencies as assets on the same chain?
There is no need to specified it's a IBC transfer？ Just by call bank.SendCoinsFromAccountToModule  or  BaseSendKeeper.SendCoins is enough? It can handle both IBC transfer and local transfer?


## Step 15: [go relayer]

### 31. 启动 relayer 节点

add the chain config files manually:
```
rly chains add --url https://raw.githubusercontent.com/cosmos/relayer/main/docs/example-configs/cosmoshub-4.json  cosmoshub
rly chains add --url https://raw.githubusercontent.com/cosmos/relayer/main/docs/example-configs/osmosis-1.json osmosis
```

生成key:
```
$ rly keys add cosmoshub key_for_cosmoshub

{"mnemonic":"much shock monitor lounge guide tribe cereal elephant slight lawn crystal robust spot dynamic dose sentence theme company drop vintage vault mother print bag","address":"cosmos1zfgr4d9sk2rt7fsl095dnnfljpm9zv7eps4xns"}


$ rly keys add osmosis key_for_osmosis
{"mnemonic":"supreme scale legal govern tag exchange vivid physical staff attract pepper gift faculty treat bike acid execute bachelor oak ozone reason follow judge vessel","address":"osmo16fzufh7qatdfpfxdctg9800ala3l85aw3tvnn5"}
```

卡在这一步，好像是拿不到 ChainProvider
```
 $ rly paths fetch
 panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x8 pc=0x4e26749]

goroutine 1 [running]:
github.com/cosmos/relayer/v2/relayer.(*Chain).ChainID(...)
	/Users/zhonglei/Codes/courses/cosmos/academy-course/relayer/relayer/chain.go:91
github.com/cosmos/relayer/v2/cmd.pathsFetchCmd.func1(0xc001052500, {0x5ec3210?, 0x0?, 0x0?})
	/Users/zhonglei/Codes/courses/cosmos/academy-course/relayer/cmd/paths.go:336 +0x829
github.com/spf13/cobra.(*Command).execute(0xc001052500, {0x5ec3210, 0x0, 0x0})
	/Users/zhonglei/go/pkg/mod/github.com/spf13/cobra@v1.5.0/command.go:872 +0x694
github.com/spf13/cobra.(*Command).ExecuteC(0xc001035400)
```

## Step 16: [cosmos js objects]

### 32. 编译前端 pb 文件

```
git submodule add git@github.com:cosmos/academy-checkers-ui.git client

$ cd scripts
$ npm init
$ npm install ts-proto@1.121.6 --save-dev --save-exact
# 之交给 makefile 完成
```


## Step 17: [enable ibc]

### 33. add an ibc module

```
ignite scaffold module leaderboard --ibc
```


### 34. create playerInfo object
```
ignite scaffold map playerInfo wonCount:uint lostCount:uint dateUpdated:string --module leaderboard --no-message
```

### 35. create board object

```
ignite scaffold single board PlayerInfo:PlayerInfo --module leaderboard --no-message
```


### 36. adjust proto file to support null filed
```
ignite generate proto-go
```


### 37. pass leaderboard keeper to checkers keeper


### 38. checkers module call the keeper of the leaderboard module

```
ignite scaffold message updateBoard --module leaderboard
```


### 39. foward player information via IBC

```
ignite scaffold packet candidate PlayerInfo:PlayerInfo --module leaderboard
```

这个命令会调整两个文件：
- tx.proto 中生成一个发送的接口（Msg类型），用于处理rpc接口
- packet.proto 中生成一个 packet 类型的对象，用于链之间的调用



### 40. IBC module intergation

方案一：
```
// app.go


// OPTIONAL: add scoped keepers in case the middleware wishes to
// send a packet or acknowledgment without
// the involvement of the underlying application	

scopedKeeperTransfer := capabilityKeeper.NewScopedKeeper("transfer")
scopedKeeperCustom1 := capabilityKeeper.NewScopedKeeper("custom1")
scopedKeeperCustom2 := capabilityKeeper.NewScopedKeeper("custom2")


// For example, if the middleware mw1 needs the ability to send a packet on custom2's port without 
// involving the underlying application custom2, it would require 
// access to the latter's scopedKeeper:
// mw1Keeper := mw1.NewKeeper(storeKey1, scopedKeeperCustom2)


// create a keeper for the stateful middleware
mw1Keeper := mw1.NewKeeper(storeKey1)
mw3Keeper := mw3.NewKeeper(storeKey3)


// instantiate the middleware as IBCModules
mw1IBCModule := mw1.NewIBCModule(mw1Keeper)
mw2IBCModule := mw2.NewIBCModule()  // optional: middleware2 is stateless middleware
mw3IBCModule := mw3.NewIBCModule(mw3Keeper)



// register the middleware in app module
// if the module maintains an independent state and/or processes sdk.Msgs
app.moduleManager = module.NewManager(
    ...
    mw1.NewAppModule(mw1Keeper),
    mw3.NewAppModule(mw3Keeper),
    transfer.NewAppModule(transferKeeper),
    custom.NewAppModule(customKeeper)
)

```


方案二：
```
// initialize base IBC applications
//
// if you want to create two different stacks with the same base application,
// they must be given different scopedKeepers and assigned different ports
transferIBCModule := transfer.NewIBCModule(transferKeeper)
customIBCModule1 := custom.NewIBCModule(customKeeper, "portCustom1")
customIBCModule2 := custom.NewIBCModule(customKeeper, "portCustom2")



stack1 := mw1.NewIBCMiddleware(mw3.NewIBCMiddleware(transferIBCModule, mw3Keeper), mw1Keeper)
// stack 2 contains mw3 -> mw2 -> custom1
stack2 := mw3.NewIBCMiddleware(mw2.NewIBCMiddleware(customIBCModule1), mw3Keeper)
// stack 3 contains mw2 -> mw1 -> custom2
stack3 := mw2.NewIBCMiddleware(mw1.NewIBCMiddleware(customIBCModule2, mw1Keeper))

ibcRouter := porttypes.NewRouter()
ibcRouter.AddRoute("transfer", stack1)
ibcRouter.AddRoute("custom1", stack2)
ibcRouter.AddRoute("custom2", stack3)
app.IBCKeeper.SetRouter(ibcRouter)

```

## [simulate production in docker]

### 41. outlines
Three independent parties - Alice, Bob, and Carol.
Two independent validator nodes, run by Alice and Bob respectively, that can only communicate with their own sentries and do not expose RPC endpoints.
Additionally, Alice's validator node uses Tendermint Key Management System (TMKMS) on a separate machine.
The two sentry nodes, run by Alice and Bob, expose endpoints to the world.
A regular node, run by Carol, that can communicate with only the sentries and exposes endpoints for use by clients.

containers:
- Alice's containers: sentry-alice, val-alice, and kms-alice.
- Bob's containers: sentry-bob and val-bob.
- Carol's containers: node-carol.

network:
- Alice's validator and key management system (KMS) are on their private network: name it net-alice-kms.
- Alice's validator and sentry are on their private network: name it net-alice.
- Bob's validator and sentry are on their private network: name it net-bob.
- There is a public network on which both sentries and Carol's node run: name it net-public.

### 42. prepare images for checkersd


### 43. prepare iamges for KMS

由于本地已经有了 rust 环境，所以不采用教程里的 multi-stage docker build 方案了，采用本地编译。
```
git clone --branch v0.12.2 https://github.com/iqlusioninc/tmkms.git
# git fetch --tags
# git checkout v0.12.2

export RUSTFLAGS=-Ctarget-feature=+aes,+ssse3
export CC_x86_64-unknown-linux-gnu=x86_64-linux-musl-gcc

# 要为交叉编译做准备，安装编译器
$ brew install SergioBenitez/osxct/x86_64-unknown-linux-gnu
# rust 工具链要安装对应的 Target
rustup target add x86_64-unknown-linux-gnu
# 修改 ~/.cargo/config，添加交叉编译的指令
[target.x86_64-unknown-linux-gnu]
linker = "x86_64-unknown-linux-gnu-gcc"

cargo build --target x86_64-unknown-linux-gnu --release --features=softsign 

生成的文件在：
/Users/zhonglei/Codes/courses/cosmos/academy-course/tmkms/target/x86_64-unknown-linux-gnu/release/tmkms

mv /Users/zhonglei/Codes/courses/cosmos/academy-course/tmkms/target/x86_64-unknown-linux-gnu/release/tmkms /Users/zhonglei/Codes/courses/cosmos/academy-course/checkers/build
```

上面的做法在 apline 的镜像下有点问题，因为 gnu 编译出来的二进制在 alpine 镜像中会有动态库缺失的问题导致无法运行，需要用 musl 这个版本来编译：
```
# ~/.cargo/config 内容改为
[build]
target = "x86_64-unknown-linux-musl"

[target.x86_64-unknown-linux-musl]
linker = "/usr/local/bin/x86_64-linux-musl-gcc"

cargo build --target x86_64-unknown-linux-musl --release --features=softsign 

mv /Users/zhonglei/Codes/courses/cosmos/academy-course/tmkms/target/x86_64-unknown-linux-musl/release/tmkms /Users/zhonglei/Codes/courses/cosmos/academy-course/checkers/build

```
更新 makefile，编译 KMS image：
```
docker-build-kms:
	docker build -f dockerfiles/prod-kms  -t tmkms_i:v0.12.2


```

运行：
```
docker run --rm -it tmkms_i:v0.12.2
```

### 44. 准备初始化的配置文件

生成配置文件：
```
echo -e node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/docker/{}:/root/.checkers \
    checkersd_i \
    init checkers
```


in genesis.json, the default initialization sets the base token to stake, so to get it to be upawn
```
docker run --rm -it \
    -v $(pwd)/docker/val-alice:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -i 's/"stake"/"upawn"/g' /root/.checkers/config/genesis.json

# macos 上要用这个命令
sed -i "" 's/"stake"/"upawn"/g' ./docker/val-alice/config/genesis.json

```

in app.toml, replace stake to upawn
```
echo -e node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/docker/{}:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei 's/([0-9]+)stake/\1upawn/g' /root/.checkers/config/app.toml

# macos 上要用这个命令
echo -e node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    sed  -Ei  "" 's/([0-9]+)stake/\1upawn/g' ./docker/{}/config/app.toml
```


Make sure that config/client.toml mentions checkers-1
```
echo -e node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/docker/{}:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei 's/^chain-id = .*$/chain-id = "checkers-1"/g' \
    /root/.checkers/config/client.toml

# macos 上用这个命令
echo -e node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    sed  -Ei "" 's/^chain-id = .*$/chain-id = "checkers-1"/g' \
    ./docker/{}/config/client.toml
```


## 45. gen Validator operator keys for Alice and Bob

create operation keys for alice, (passpharse: password)

```
docker run --rm -it \
    -v $(pwd)/docker/val-alice:/root/.checkers \
    checkersd_i \
    keys \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    add alice
Enter keyring passphrase:
Re-enter keyring passphrase:

- name: alice
  type: local
  address: cosmos1dngyd53yxc5c9vswss9a3nalrpqdwfzqj7k2e5
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A8oMqt3oVk6R53WW9UaZ8+oE7S9q3fqZZcaTxFalTMBz"}'
  mnemonic: ""


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

busy sugar tobacco dizzy spray coin among salute stadium slot festival walk unfair bench bargain river tribe birth pond solar royal plug fatigue priority
```

create operation keys for bob, (passpharse: password)
```
$ mkdir -p docker/val-bob/keys
$ echo -n password > docker/val-bob/keys/passphrase.txt
$ docker run --rm -it \
    -v $(pwd)/docker/val-bob:/root/.checkers \
    checkersd_i \
    keys \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    add bob

Enter keyring passphrase:
Re-enter keyring passphrase:

- name: bob
  type: local
  address: cosmos1u2yzjj59v7qudl4mgk3svtz0hjxptrvz0zpqzu
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A0zXFTPghu0VKz2GsbYAqA6N8ihF5FNZlID+ivGBzxvh"}'
  mnemonic: ""


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

mixture arena announce spin fringe fatigue thrive tragic energy cable oak camera gorilla term elite august neutral razor client isolate cradle casino catch feed
```


## 46. 为 Alice 的 tmkms 生成配置文件

```
docker run --rm -it \
    -v $(pwd)/docker/kms-alice:/root/tmkms \
    tmkms_i:v0.12.2 \
    init /root/tmkms
   Generated KMS configuration: /root/tmkms/tmkms.toml
   Generated Secret Connection key: /root/tmkms/secrets/kms-identity.key
```

--- 
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
[play with IBC tokens]: https://interchainacademy.cosmos.network/hands-on-exercise/2-ignite-cli-adv/8-wager-denom.html
[go relayer]: https://interchainacademy.cosmos.network/hands-on-exercise/5-ibc-adv/3-go-relayer.html
[cosmos js objects]: https://interchainacademy.cosmos.network/hands-on-exercise/3-cosmjs-adv/1-cosmjs-objects.html
[enable ibc]: https://interchainacademy.cosmos.network/hands-on-exercise/5-ibc-adv/8-ibc-app-checkers.html
[simulate production in docker]: https://interchainacademy.cosmos.network/hands-on-exercise/4-run-in-prod/1-run-prod-docker.html