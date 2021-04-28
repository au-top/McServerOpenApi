# McServerOpenApi

可以在 plugin 中自定义插件 参考 OpenMods 代码用作示例
自带 OpenMods 插件用于向下同步 mods

/src/config.json

mcServerPath 对应服务器路径
``` 
{
    "mcServerPath":"D:\\project\\js\\mcServerOpenSourceFile\\test"
}
```



# Plugin

## OpenMods

在MC服务器目录中创建 

    clientMods 文件夹用于放置客户端独有Mod
    mods       文件夹放置服务器独有Mod (mods 一般服务器mod端自带)

在client/openMods 存在使用Go 编写的mod同步客户端 支持 GNU/Linux , Windows


### CLI
```
    --config //配置文件路径
    --gamePath //游戏路径
```

### 配置文件
```
client/openMods/config.json
```

### 使用示例
```
.\fromMcServerGetModConfig.exe -config "./config.json" -gamePath "D:\HMCL\.minecraft"
```

### 其他
    客户端支持Mods下载缓存 位于 ./cache 文件夹
    Mods备份 在mc游戏目录中创建更新前的mod备份 mcgame/back_{time_token}

### Api
```
    get Mod List
    GET {{host}}/modInfo 

    [
        {
            "filename": "AttributeFix-1.16.5-10.1.2.jar",
            "length": 9554
        },
        {
            "filename": "BackTools-1.16.5-10.1.0.jar",
            "length": 23874
        }
    ]


    download Mod File
    GET {{host}}/modSource
    [buffer]
```