
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

## breakpoint:
https://interchainacademy.cosmos.network/hands-on-exercise/1-ignite-cli/3-stored-game.html#some-initial-thoughts