# DDD模块开发助手使用示例

## 基础使用示例

### 示例 1：学习进度跟踪模块开发

**场景**：开发一个学习进度跟踪功能模块

```
/ddd-develop-module
```

**第一轮对话 - 业务理解：**
```
Skills: "您好！我是DDD模块开发助手。请描述您想要开发的业务模块。"
用户: "我想开发一个学习进度跟踪功能，用来记录学生的练习情况和掌握程度。"

Skills: "很好的需求！为了更好地设计，请问：
1. 主要服务对象是学生还是老师？
2. 需要跟踪哪些具体的学习活动？
3. 期望达到什么样的业务目标？"
用户: "主要服务学生，跟踪做题练习，目标是帮助学生了解自己的学习进展。"
```

**第二轮对话 - 领域概念识别：**
```
Skills: "基于您的描述，我识别出以下潜在的领域概念：
🎯 核心实体：学习会话、练习题目、学生档案
📊 值对象：掌握程度(百分比)、学习时长、练习得分
⚡ 领域事件：练习开始、题目完成、掌握度更新
📁 聚合建议：学习进度聚合(会话+统计数据)

请问这些概念是否准确？还有其他重要的业务概念吗？"
用户: "概念基本准确，但掌握程度应该是一个更复杂的计算结果。"
```

**第三轮对话 - 聚合边界确认：**
```
Skills: "关于聚合设计，我建议：
📦 学习进度聚合：包含学习会话、练习记录、统计计算
👤 学生聚合：学生基本信息、学习偏好设置
📘 题目聚合：题目信息、知识点关联

这样的边界划分是否合理？您有什么调整建议？"
用户: "学习进度聚合的范围有点大，能否拆分得更细一些？"
```

**第四轮对话 - 业务规则澄清：**
```
Skills: "重新调整后的聚合设计：
📝 会话聚合：单次学习会话的完整生命周期
📊 统计聚合：跨会话的学习数据汇总
📈 掌握度聚合：知识点掌握情况计算

关键业务规则识别：
❌ 同一学生对同一题目只能有一个活跃会话
📈 掌握度 = (正确题目数 / 总题目数) × 100%
⏰ 学习时长 = 最后答题时间 - 开始时间

这些规则是否完整准确？"
用户: "规则很准确，但还需要考虑连续答对的奖励机制。"
```

**第五轮对话 - 技术要求确认：**
```
Skills: "技术集成要求确认：
🔗 集成现有用户服务获取学生信息
📡 通过领域事件通知推荐系统
📚 生成RESTful API接口
💾 使用GORM进行数据持久化
🧪 包含完整的单元测试

这些技术要求是否符合您的预期？"
用户: "完全符合预期，开始生成代码吧！"
```

**第六轮 - 代码生成执行：**
```
Skills: "确认无误，开始生成代码..."

[执行完整的代码生成流程，创建以下文件结构]

internal/
├── domain/
│   └── learning_progress/
│       ├── aggregate/
│       │   ├── session_aggregate.go      # 会话聚合根
│       │   ├── statistics_aggregate.go   # 统计聚合根
│       │   └── proficiency_aggregate.go  # 掌握度聚合根
│       ├── entity/
│       │   ├── learning_session.go       # 学习会话实体
│       │   ├── exercise_record.go        # 练习记录实体
│       │   └── student_profile.go        # 学生档案实体
│       ├── valueobject/
│       │   ├── proficiency_score.go      # 掌握度分数
│       │   ├── study_duration.go         # 学习时长
│       │   └── exercise_result.go        # 练习结果
│       ├── event/
│       │   ├── session_started.go        # 会话开始事件
│       │   ├── exercise_completed.go     # 练习完成事件
│       │   └── proficiency_updated.go    # 掌握度更新事件
│       └── service/
│           └── proficiency_calculator.go # 掌握度计算服务
├── application/
│   └── learning_progress/
│       └── service/
│           ├── progress_tracking_service.go  # 进度跟踪应用服务
│           └── statistics_service.go         # 统计分析应用服务
├── infrastructure/
│   └── persistence/
│       └── gorm/
│           └── learning_progress/
│               ├── dao/
│               │   ├── session_dao.go        # 会话DAO接口
│               │   ├── record_dao.go         # 记录DAO接口
│               │   └── statistics_dao.go     # 统计DAO接口
│               └── model/
│                   ├── session_model.go      # 会话数据模型
│                   ├── record_model.go       # 记录数据模型
│                   └── statistics_model.go   # 统计数据模型
└── interfaces/
    └── http/
        └── learning_progress_handler.go      # HTTP接口处理器

[同时生成配套资源]
├── migrations/20260126100002_create_learning_progress_tables.sql
├── docs/api/learning_progress_swagger.json
├── internal/domain/learning_progress/test/
│   ├── aggregate/session_aggregate_test.go
│   └── service/proficiency_calculator_test.go
```

### 示例 2：用户积分系统开发

**场景**：为现有用户系统添加积分功能

```
/ddd-develop-module --module-name user_points
```

**简化版对话流程：**
```
用户: "需要给用户系统添加积分功能，支持赚取和消费积分。"

Skills: "识别到核心概念：
💰 积分账户(聚合根)
💳 积分流水(实体)
🏆 积分等级(值对象)
⚡ 积分变动事件

确认采用积分账户聚合设计？"
用户: "确认"

Skills: "生成积分系统代码..."
[执行代码生成]
```

## 高级使用示例

### 示例 3：复杂电商订单系统

**场景**：开发完整的电商订单处理模块

```
/ddd-develop-module --module-name ecommerce_order
```

**多轮深度对话：**
```
# 第一轮：业务复杂度分析
用户: "需要一个完整的电商订单系统，包含下单、支付、发货、售后等全流程。"

Skills: "识别到复杂业务场景，建议分阶段开发：
Phase 1: 订单创建和基础流程
Phase 2: 支付集成和状态管理
Phase 3: 物流跟踪和售后服务

您希望采用哪种方式？"
用户: "分阶段开发，先实现基础订单流程。"

# 第二轮：核心领域识别
Skills: "基础订单流程的核心概念：
🛒 订单聚合：订单主信息、商品项、价格计算
🚚 配送聚合：配送地址、物流信息
💳 支付聚合：支付信息、交易记录
👨‍💼 客户聚合：客户信息、收货地址

建议先聚焦订单聚合的实现？"
用户: "同意，先实现订单聚合。"

# 第三轮：业务规则细化
Skills: "订单聚合的业务规则：
❌ 商品库存充足才能下单
💰 订单金额计算包含运费
📅 订单创建后24小时内未支付自动取消
🔄 已支付订单支持取消申请

这些规则是否完整？"
用户: "还需要添加优惠券使用规则。"

# 执行生成
Skills: "开始生成基础订单系统代码..."
[生成相应的代码结构]
```

### 示例 4：与现有Skills集成

**场景**：结合database-migrator和api-doc-generator

```
# 启动DDD开发流程
/ddd-develop-module

# 在代码生成阶段自动调用其他Skills
Skills: "正在生成代码...
[调用 /database-migrator 创建表结构]
[调用 /api-doc-generator 生成API文档]
代码生成完成！"
```

## 团队协作示例

### 示例 5：多人协作的领域建模

**场景**：团队成员分别负责不同方面的建模

```
# 团队成员A：业务分析师角色
/ddd-analyze-business "在线教育平台的作业提交功能"

# 团队成员B：技术架构师角色
/ddd-design-domain --entities "作业、学生、教师" --concepts "提交、批改、反馈"

# 团队负责人：整合和确认
/ddd-export-model --format json
/ddd-validate-design --strict

# 最终执行代码生成
/ddd-generate-code --from-session combined_analysis
```

### 示例 6：DDD学习和实践

**场景**：新手学习DDD概念

```
# 学习模式启动
/ddd-develop-module --learning-mode

Skills: "欢迎来到DDD学习模式！我们将通过一个简单的例子来理解DDD概念。

请描述一个您熟悉的简单业务场景..."
用户: "图书馆借书系统"

Skills: "很好！让我们一步步分析：
第一步：识别核心业务概念
您认为借书系统中最重要的是什么？"
[通过引导式提问帮助用户理解DDD概念]
```

## 故障排除示例

### 示例 7：处理设计争议

**场景**：团队对聚合边界有分歧

```
用户: "我们认为订单和支付应该放在同一个聚合中。"

Skills: "检测到聚合设计争议，让我帮您分析：

当前建议：订单聚合 + 支付聚合（分离设计）
您的观点：订单支付聚合（合并设计）

分析对比：
✅ 分离优势：职责清晰、可独立扩展、降低复杂度
⚠️ 合并风险：聚合过大、事务边界模糊、耦合度高

建议进行小范围原型验证，您觉得如何？"
```

### 示例 8：代码生成问题处理

**场景**：生成的代码不符合预期

```
用户: "生成的实体缺少验证逻辑。"

Skills: "检测到代码质量问题，提供以下解决方案：

1. 重新配置生成选项：
   /ddd-generate-code --with-validation --force

2. 手动添加验证逻辑：
   [提供验证代码模板]

3. 更新全局配置：
   [指导修改config.yaml]

您希望采用哪种方式？"
```

## 性能优化示例

### 示例 9：大批量代码生成

**场景**：同时开发多个相关模块

```
# 批量模式启动
/ddd-develop-module --batch-mode

Skills: "检测到批量开发需求，请提供模块列表：
1. 用户管理模块
2. 内容管理模块
3. 评论系统模块

将为每个模块创建独立的会话并行处理..."
[并行执行多个模块的开发流程]
```

### 示例 10：增量更新现有模块

**场景**：为已有DDD模块添加新功能

```
# 增量开发模式
/ddd-develop-module --existing-module user_management --incremental

Skills: "检测到现有模块 user_management，当前结构：
- UserAggregate (用户聚合)
- UserProfile (用户档案实体)
- UserStatus (用户状态值对象)

请描述要添加的新功能..."
用户: "添加用户等级和积分系统"

Skills: "识别增量变更：
➕ 新增：UserLevel(值对象)、Points(实体)
🔄 修改：UserAggregate(添加积分相关行为)
⚡ 新增：PointsEarnedEvent(领域事件)

确认变更方案后开始增量生成..."
```

这些示例展示了DDD模块开发助手在各种实际场景中的灵活应用，从简单功能到复杂系统，从个人开发到团队协作，都能提供专业的DDD开发支持。