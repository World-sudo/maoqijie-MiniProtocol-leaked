思考和回复始终使用中文
所有文件不得大于300行（需求文档除外），要模块化进行拆分
文件不能使用下划线命名，作为模块的话要单独开目录

# 代码拆分与命名（重要）

- 不要用“前后缀堆叠”来组织代码（例如 `subsitesReleaseRolloutsAdminHandlers.go` 这类命名）；当一个功能需要反复加业务前缀/后缀才能描述清楚时，必须拆分为独立模块目录。
- 模块化优先但避免“目录爆炸”：能在同一目录用多个短文件表达清楚的，优先直接 `xx.go`（例如 `internal/api/auth.go`、`internal/api/jobs.go`）。只有当某个功能域需要多文件/子域明显/需要复用时，才单独开目录（例如 `internal/api/subsites/`、`frontend/src/pages/SubsitesAdminPage/`），目录内文件名保持简洁通用（例如 `handlers.go`、`routes.go`、`models.go`、`service.go`、`types.ts`、`index.ts`）。
- 单文件接近 300 行时必须继续拆分，避免“大杂烩文件”；同一目录下文件数量明显变多时，优先再按子域拆分目录，保持可读性与可维护性。

## 前缀命名治理（强制）

- 在同一目录内，**出现 3 个及以上同业务前缀文件**（如 `releaseCreate.go`、`releaseSource.go`、`releaseGitHub.go`）时，必须立即拆为子模块目录（如 `release/`），并改为通用短名（如 `create.go`、`source.go`、`github.go`、`versions.go`）。
- **禁止继续新增** `releaseXxx.go`、`deployXxx.go`、`subsiteXxx.go` 这类“业务前缀 + 行为后缀”堆叠文件；已有 2 个同前缀文件时，新增第 3 个前必须先完成目录化拆分。
- 以上规则同时适用于前端（`*.ts`/`*.tsx`）：例如出现 `releaseSourceCard.tsx`、`releaseCreateCard.tsx`、`releaseListTable.tsx` 时，应拆到 `release/` 目录并使用语义化短名。
- 若为紧急修复允许临时保留旧命名，必须在同次任务内补拆分；除非用户明确要求“本次不要重构目录结构”。
- 反例：`internal/api/subsites/releaseCreate.go`、`internal/api/subsites/releaseSourceGitHubList.go`。
- 正例：`internal/api/subsites/release/create.go`、`internal/api/subsites/release/source.go`、`internal/api/subsites/release/versions.go`。