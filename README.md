### 目录结构

```
|api     外部传输结构体定义如：proto文件等
|app     组件启动模块，包含CLI相关操作
|pkg     
|---config       配置目录 
|   |----app     对应模块配置
|   |----common  通用配置
|---core         核心业务逻辑
|   |----shunt   
|        |---adapter     接口防腐层
|            |---facade  对接外部服务
|            |---http    web层
|                |--handlers 控制层
|                |--routers  路由
|                |--server   服务入口 
|            application 业务层
|                    cqe
|                    |---cmd    新增、修改结构体
|                    |---query  查询结构体
|                    dto        内部传输结构体
|                    service    业务服务实现
|            repository  持久层
|                    entity  实体
|                    xxxRepo 持久化接口
|---log      日志框架
|---plugin   三方依赖
|---types    常量、自定义类型
|---util     工具类     
                    
```
