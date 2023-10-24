### 目录结构

```
|api     外部传输结构体定义如：proto文件等
|app     
|   |----shunt
|        |---cmd     组件启动模块，包含CLI相关操作
|        |---config     配置文件及运行时配置
|        |---adapter     接口防腐层
|            |---cqe
|                |---cmd    新增、修改结构体
|                |---query  查询结构体
|            |---facade  对接外部服务
|            |---http    web层
|                |--handlers 控制层
|                |--middleware 中间件
|                |--routers  路由
|                |--server   服务入口 
|                    service    业务服务实现
|        |---repository  持久层
|                  entity  实体
|                  xxxRepo 持久化接口
|        |---service  功能层
|                |--- bus  业务逻辑
|                   |--- impl 具体实现
|                |--- dto  传输DTO

|pkg     
|---config   配置加载工具类 
|---plugin   三方插件
|   |---es
|   |---log
|   |---mongo
|   |---mysql
|   |---redis
|---runtime   运行时插件
|   |---component 组件
|   |---events 事件总线
|   |---watchdog 定时任务
|---types    常量、自定义类型
|---util     工具类     
|   |---cache 内存缓存
|   |---cache 文件处理
|   |---http http工具包
|   |---jwt jwt工具包
|   |---snowflake 雪花算法UUID
|   |---xxxx 其他工具类                
```
