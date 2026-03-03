# MathFun API 接口设计规范

**文档目的:**
本文档旨在为 MathFun 项目的 API 接口设计提供统一的标准和指导，确保接口的规范性、一致性、可读性和可维护性。

**版本:** 1.0  
**生效日期:** 2026-02-02

## 1. 术语定义
*   **API:** Application Programming Interface，应用程序编程接口。
*   **Client:** API 的调用方，例如前端应用、移动端应用或其他后端服务。
*   **Server:** API 的提供方，即后端服务。
*   **Endpoint:** 一个具体的 API 地址，例如 `/api/knowledge/domains`。
*   **DTO (Data Transfer Object):** 用于接口间数据传输的对象。

## 2. 通用规范

*   **2.1 风格:** 采用 RESTful 设计风格。
*   **2.2 协议:** 使用 HTTPS 协议。
*   **2.3 命名:**
    *   URL 路径使用小写字母和连字符 (`kebab-case`)，例如 `/api/knowledge/knowledge-nodes`。
    *   **JSON 请求/响应体中的字段名使用小驼峰命名法 (`camelCase`)，例如 `nodeName`, `createdAt`。**
    *   **URL 查询参数名和路径参数名也使用小驼峰命名法 (`camelCase`)，例如 `?parentId=123&includeDetails=true`，`/api/knowledge/trunks/{domainId}`。** (新增)
*   **2.4 版本控制:** API 版本号应体现在 URL 路径中，格式为 `/api/v{version_number}/...`，例如 `/api/v1/knowledge/domains`。当前默认版本为 `v1`。
*   **2.5 字符编码:** 所有请求和响应均使用 UTF-8 编码。

## 3. 请求规范

*   **3.1 HTTP 方法:**
    *   `GET`: 获取资源。
    *   `POST`: 创建资源或执行非幂等的操作（如提交表单、触发动作）。
    *   `PUT`: 全量更新资源或创建（若资源不存在）。
    *   `PATCH`: 部分更新资源。
    *   `DELETE`: 删除资源。
*   **3.2 请求头 (Headers):**
    *   `Content-Type`: 指定请求体的数据格式，如 `application/json`。
    *   `Authorization`: 用于携带认证凭据（如 Bearer Token）。
    *   `X-Request-ID` (可选): 客户端生成的请求唯一标识，用于链路追踪。
*   **3.3 请求体 (Body):**
    *   通常使用 `application/json` 格式。
    *   `GET` 和 `DELETE` 请求一般不带请求体（除非有特殊需求，如批量删除）。
*   **3.4 查询参数 (Query Parameters):**
    *   **分页参数:**
        *   `page` (integer, optional, default: 1): 页码，从 1 开始。
        *   `pageSize` (integer, optional, default: 20, max: 100): 每页大小。
    *   **排序参数:**
        *   **参数名:** `sort`
        *   **参数值:** 采用 `,` 分隔的字段列表，字段前可加 `-` 表示降序，无符号表示升序。
        *   **示例:**
            *   单一字段降序: `?sort=-createdAt`
            *   单一字段升序: `?sort=createdAt` (或 `?sort=+createdAt`，`+` 可省略)
            *   多字段排序 (先按 priority 降序，再按 name 升序): `?sort=-priority,name`
    *   **过滤参数:**
        *   **参数名:** `filter[字段名]`
        *   **参数值:**
            *   **单值等效过滤:** `值`。例如 `?filter[name]=john`。
            *   **多值 IN 过滤:** 使用 `,` 分隔多个值，表示查询该字段的值在列表中的记录。例如 `?filter[category]=math,science` (category 为 math 或 science)。
            *   **范围/比较/高级操作:** 服务端应支持通过在 `filter[字段名]` 的值前添加特定前缀来表示不同的比较操作。
                *   **推荐前缀 (参考 laravel-query-builder):**
                    *   `>` 或 `gt`: 大于。例如 `?filter[salary]=gt:3000` (薪资大于 3000)。
                    *   `<` 或 `lt`: 小于。例如 `?filter[age]=lt:18` (年龄小于 18)。
                    *   `>=` 或 `gte`: 大于等于。例如 `?filter[salary]=gte:3000` (薪资大于等于 3000)。
                    *   `<=` 或 `lte`: 小于等于。例如 `?filter[age]=lte:65` (年龄小于等于 65)。
                    *   `!=` 或 `not`: 不等于。例如 `?filter[status]=not:inactive` (状态不是 inactive)。
                    *   `like`: 模糊匹配 (注意：实现时需防止 SQL 注入)。例如 `?filter[name]=like:%john%` (名字包含 john)。
                    *   `in`: 包含在列表中 (等同于多值 OR，但语义更明确)。例如 `?filter[id]=in:1,2,3` (id 为 1 或 2 或 3)。 (注意：此功能与 `filter[字段名]=值1,值2,值3` 等价，但 `in` 前缀语义更明确)
                    *   `between`: 在某个范围内 (包含边界)。例如 `?filter[salary]=between:3000,5000` (薪资在 3000 到 5000 之间)。
                    *   `date:`: 日期比较 (精确到天)。例如 `?filter[createdAt]=date:2023-10-26` (创建日期为 2023-10-26)。
                    *   `date>:date<:, date>=:, date<=:`: 日期范围比较。例如 `?filter[createdAt]=date>=:2023-01-01&filter[createdAt]=date<:2024-01-01` (创建日期在 2023 年)。
        *   **逻辑组合:**
            *   **AND:** 不同的 `filter[字段名]` 参数之间隐式构成 AND 逻辑。例如 `?filter[status]=active&filter[age]=gte:18` (状态为 active 且年龄 >= 18)。
            *   **OR (IN):** 同一 `filter[字段名]` 的值使用逗号分隔表示 OR 逻辑（即 `IN` 查询）。例如 `?filter[category]=math,science` (category 为 math 或 science)。
        *   **注意事项:**
            *   **安全:** 服务端必须对所有过滤字段和操作符进行严格的白名单校验和输入净化，防止 SQL 注入、XSS 等攻击。
            *   **性能:** 复杂的过滤查询（特别是模糊匹配、多字段联合查询）可能对性能有较大影响，需要进行评估和优化，必要时添加数据库索引。
            *   **实现:** 后端解析时，应根据 `filter[字段名]` 的值中是否存在前缀（如 `gt:`, `between:`, `like:` 等）来决定构建何种类型的数据库查询条件（如 `>`, `BETWEEN`, `LIKE` 等）。
    *   **其他参数:** 根据具体业务需求定义。

## 4. 响应规范

*   **4.1 响应头 (Headers):**
    *   `Content-Type`: 指定响应体的数据格式，如 `application/json`。
*   **4.2 响应体 (Body):**
    *   **成功响应:**
        *   `GET /resource`: 返回资源对象或资源对象数组。
        *   `POST /resource`: 返回创建成功的资源对象。
        *   `PUT /resource/{id}`, `PATCH /resource/{id}`: 返回更新后的资源对象。
        *   `DELETE /resource/{id}`: 返回标准成功格式（见下文）。
        *   **标准成功格式 (Standard Success Response):**
            ```json
            {
              "success": true,
              "code": 0, // 业务状态码，0表示成功
              "message": "Success message", // 可选的成功提示信息
              "data": { /* 实际返回的业务数据 */ }, // 对于分页列表，data结构为 { "list": [...], "pagination": { "page": 1, "pageSize": 20, "total": 100 } }
              "requestId": "..." // 可选的请求ID，用于链路追踪
            }
            ```
    *   **失败响应:**
        *   **标准失败格式 (Standard Error Response):**
            ```json
            {
              "success": false,
              "code": 40001, // 业务错误码
              "message": "Error message", // 错误描述信息
              "requestId": "..." // 可选的请求ID
            }
            ```
        *   `data` 字段在失败时通常不包含业务数据，可根据需要包含错误详情等信息。

## 5. HTTP 状态码

*   `200 OK`: 请求成功。
*   `201 Created`: 资源创建成功。
*   `204 No Content`: 请求处理成功，但无返回内容（如 DELETE 操作）。
*   `400 Bad Request`: 请求参数错误或格式不正确。
*   `401 Unauthorized`: 未认证或认证失败。
*   `403 Forbidden`: 已认证但无权访问。
*   `404 Not Found`: 请求的资源不存在。
*   `422 Unprocessable Entity`: 请求格式正确，但语义错误（如校验失败）。
*   `429 Too Many Requests`: 请求频率超过限制。
*   `500 Internal Server Error`: 服务器内部错误。

## 6. 错误处理

*   **6.1 错误码:** 使用统一的错误码体系，区分系统错误和业务错误。
    *   系统错误码 (5xx): 通常由 HTTP 状态码表示，或使用 `50000` 起始的业务码。
    *   业务错误码 (4xx): 使用 `40000` 起始的三位或四位数字，例如 `40001` (参数错误), `40100` (用户未找到)。
*   **6.2 错误信息:** 提供清晰、友好的错误信息，便于 Client 端理解和处理。避免暴露敏感的系统信息。

## 7. 安全规范

*   **7.1 认证:** 所有需要认证的接口必须通过 `Authorization` Header 携带有效的 Token。
*   **7.2 授权:** 根据用户角色和权限控制对资源的访问。
*   **7.3 输入校验:** 对所有来自 Client 的输入（URL、Query、Body、Header）进行严格校验，防止 SQL 注入、XSS 等攻击。
*   **7.4 敏感数据:** 对敏感数据（如密码、身份证号）进行脱敏处理，不在日志中明文打印。

## 8. 性能与监控 (可观测性)

*   **8.1 速率限制:** 对写操作和可能消耗大量资源的读操作实施速率限制。
*   **8.2 链路追踪:** 建议在请求头中传递 `X-Request-ID` 并在服务端日志中记录，便于问题排查。
*   **8.3 日志:** 记录关键的请求、响应和错误日志。**必需对关键业务字段（如入参、出参、执行耗时）进行打点。**

## 9. 兼容性与版本管理

*   **9.1 版本管理:** API 版本号体现在 URL 中。进行不兼容更新时，应发布新版本 (如 v2)，旧版本需提供**灰度下线流程**（通知 → 只读 → 禁用 → 清理）。
*   **9.2 向后兼容:** 尽可能保持向后兼容，避免破坏性变更。

## 10. 文档

*   **10.1 生成:** 使用 Swagger/OpenAPI 规范自动生成 API 文档。
*   **10.2 维护:** 保持代码注释与 API 文档的同步更新。