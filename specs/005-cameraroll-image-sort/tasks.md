# Actionable Tasks: 全画像を対象とした日付ソート機能

**Branch**: `005-cameraroll-image-sort`
**Feature**: 全画像を対象とした日付ソート機能

## Phase 1: Setup

- [X] T001 既存のコードベースを理解するため、`src/viewer/cmd/viewer/cameraroll.go` と `src/viewer/internal/model/thumbnail.go` の内容を確認する

## Phase 2: User Story 1 - 全画像のシームレスな日付順表示

**Goal**: ユーザーがカメラロールを開いた際に、すべての画像が日付の新しい順にソートされて表示され、スクロールすることで過去の画像をシームレスに閲覧できる。
**Independent Test**: `GET /cameraroll` エンドポイントを `last_evaluated_key` パラメータ付きで連続して呼び出し、返される `thumbnails` リストが常にタイムスタンプの降順でソートされており、重複がないことを確認する。

### Implementation Tasks

- [X] T002 [US1] `src/viewer/internal/model/thumbnail.go` の `ListThumbnails` 関数を修正し、DynamoDBから取得した結果をタイムスタンプの **降順** でソートするように変更する
- [X] T003 [US1] `src/viewer/cmd/viewer/cameraroll.go` の `CameraRollHandler` を修正し、`last_evaluated_key` と `limit` クエリパラメータをリクエストから受け取れるようにする
- [X] T004 [US1] `src/viewer/cmd/viewer/cameraroll.go` の `CameraRollHandler` 内で、`model.ListThumbnails` を呼び出す際に、受け取った `last_evaluated_key` と `limit` を渡すように修正する
- [X] T005 [US1] `src/viewer/cmd/viewer/cameraroll.go` の `CameraRollHandler` のレスポンスに、`ListThumbnails` から返された次のページのキー `last_evaluated_key` を含めるように修正する

### Testing Tasks (TDD)

- [ ] T006 [P] [US1] `src/viewer/internal/model/thumbnail_test.go` に、`ListThumbnails` が結果を降順で返すことを検証するユニットテストを追加する
- [ ] T007 [P] [US1] `src/viewer/cmd/viewer/cameraroll_test.go` に、`CameraRollHandler` がページネーションのパラメータを正しく処理し、レスポンスに `last_evaluated_key` を含めることを検証するユニットテストを追加する

## Phase 3: Polish & Cross-Cutting Concerns

- [X] T008 `serverless.yml` を確認し、`GET /cameraroll` エンドポイントの定義が正しいことを確認する
- [ ] T009 `serverless deploy` を実行して変更をAWSにデプロイし、実際にエンドポイントを叩いて動作確認を行う

## Dependencies

- **User Story 1** は自己完結しており、他のユーザーストーリーへの依存はありません。

## Parallel Execution

- **User Story 1 内**:
  - `T006` (モデル層のテスト) と `T007` (ハンドラ層のテスト) は、それぞれの実装タスク（`T002`, `T003-T005`）と並行して、または先んじて（TDDスタイル）進めることができます。

## Implementation Strategy

- **MVP First**: User Story 1 がこの機能のMVP（Minimum Viable Product）です。Phase 2 のタスクを完了させることで、主要な価値を提供できます。
- **Incremental Delivery**: まずはバックエンドのAPI改修を完了させ、その後フロントエンド（このドキュメントのスコープ外）の無限スクロール実装へと進むことができます。
