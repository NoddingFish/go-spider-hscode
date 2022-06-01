# go-spider-hscode

爬取 [HS编码查询](https://www.hsbianma.com/) 网站

> ## 注意
>
> 1. 首页是 JS 动态加载，是直接获取 HSStaticData.js
> 2. 爬取过程中请注意休眠时间，尽量不要影响到他人



```shell
go run main.go
```



运行过程中会生成两个 `sql` 文件，分别是 `hs_code1.sql` 和 `hs_code_list1.sql`
