# Actionable Tasks: カメラロールのページネーション修正

**Feature**: [カメラロールのページネーション修正](spec.md)
**Plan**: [Implementation Plan](plan.md)

---

## Phase 1: Foundational Tasks (前提タスク)

- [ ] T001 既存のページネーションの問題を再現する単体テストを作成する in `src/viewer/cmd/viewer/cameraroll_test.go`

## Phase 2: User Story 1 - 正しい時系列での写真表示 (US1)

**Goal**: ユーザーがページネーションを利用して、常に時系列が正しい順序で写真を表示できるようにする。
**Independent Test**: 「次へ」ボタンをクリックするたびに、表示される写真が常に現在表示されている写真よりも古い日付のものだけで構成されていることを確認する。

### Implementation Tasks

- [ ] T002 [US1] DynamoDBのソート順を利用するようにデータ取得ロジックを修正する in `src/viewer/internal/model/thumbnail.go`
- [ ] T003 [US1] カーソルベースのページネーションをハンドラに実装する in `src/viewer/cmd/viewer/cameraroll.go`
- [ ] T004 [US1] `next_cursor` を使って「次へ」ボタンのリンクを生成するようにテンプレートを修正する in `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`

## Phase 3: Polish & Verification (仕上げと検証)

- [ ] T005 T001で作成したテストが成功することを確認する in `src/viewer/cmd/viewer/cameraroll_test.go`
- [ ] T006 [quickstart.md](quickstart.md) の手順に従い、ブラウザで手動テストを実施する

---

## Dependencies & Execution Strategy

### Story Completion Order

1.  **User Story 1 (US1)**: このフィーチャーの唯一のユーザーストーリーであり、すべてのタスクが含まれます。

### Implementation Strategy

この修正は単一のユーザーストーリーに集約されているため、タスクはT001からT006まで順番に実行されます。TDD（テスト駆動開発）のアプローチに基づき、まず失敗するテスト（T001）を作成し、次に実装（T002-T004）を行い、最後にテストが成功することを確認（T005-T006）します。

### Parallel Execution

このフィーチャーのタスクは、バックエンドのロジックとフロントエンドのテンプレートが密接に関連しているため、並行実行は推奨されません。記載された順序での実行が最も効率的です。
