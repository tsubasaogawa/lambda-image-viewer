# Tasks: Camera Roll デザイン改修

**Input**: Design documents from `/specs/004-cameraroll-design-update/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, quickstart.md

**Tests**: TDD（テスト駆動開発）アプローチに基づき、実装前にテストタスクを含めます。

**Organization**: タスクはユーザーストーリーごとにグループ化され、各ストーリーの独立した実装とテストを可能にします。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 並行して実行可能（異なるファイル、未完了のタスクへの依存なし）
- **[Story]**: このタスクが属するユーザーストーリー（例：US1, US2）
- 説明には正確なファイルパスを含める

## Path Conventions

- `src/viewer/` と `tests/viewer` をリポジトリルートからの相対パスとして使用します。

---

## Phase 1: Setup (共有インフラ)

**Purpose**: 既存コードの理解と準備

- [ ] T001 既存の `cameraroll.go` と `camera_roll.html.tmpl` のコードをレビューし、現在の実装を理解する

---

## Phase 2: User Story 1 - 月単位での写真閲覧 (Priority: P1) 🎯 MVP

**Goal**: バックエンドで写真データを月ごとにグループ化し、フロントエンドに渡すためのデータ構造を準備する。

**Independent Test**: `cameraroll_test.go` の単体テストを実行し、写真が正しく月ごとにグループ化され、ソートされることを確認する。

### Tests for User Story 1 (TDD) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T002 [P] [US1] `tests/viewer/cameraroll_test.go` に、複数月にまたがる写真が正しく年月でグループ化されることを検証するテストケースを追加する
- [ ] T003 [P] [US1] `tests/viewer/cameraroll_test.go` に、撮影日がない写真が「撮影日不明」グループに分類されることを検証するテストケースを追加する
- [ ] T004 [P] [US1] `tests/viewer/cameraroll_test.go` に、写真グループが新しい順（降順）にソートされることを検証するテストケースを追加する

### Implementation for User Story 1

- [ ] T005 [US1] `src/viewer/cmd/viewer/cameraroll.go` を修正し、DynamoDBから取得したすべての写真アイテムを月ごとにグループ化するロジックを実装する (T002に依存)
- [ ] T006 [US1] `src/viewer/cmd/viewer/cameraroll.go` に、EXIFデータがない写真を「撮影日不明」グループとして処理するロジックを追加する (T003に依存)
- [ ] T007 [US1] `src/viewer/cmd/viewer/cameraroll.go` で、グループ化した写真データを新しい順にソートするロジックを実装する (T004に依存)
- [ ] T008 [US1] `src/viewer/cmd/viewer/cameraroll.go` で、最終的なデータ構造をHTMLテンプレートに渡すように修正する

**Checkpoint**: この時点で、User Story 1のバックエンドロジックは完全に機能し、単体テストがすべて成功するはずです。

---

## Phase 3: User Story 2 - グリッドレイアウトでの写真表示 (Priority: P1) 🎯 MVP

**Goal**: User Story 1で準備されたデータ構造を使用し、Flickrのようなレスポンシブなグリッドレイアウトで写真を表示する。

**Independent Test**: Lambda関数をデプロイし、ブラウザでカメラロールページにアクセスする。写真が月ごとにグループ化され、ヘッダーが表示され、画面サイズに応じてグリッドが正しく変化することを目視で確認する。

### Implementation for User Story 2

- [ ] T009 [US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の既存のレイアウトを削除し、新しいグリッドレイアウトの骨格となるHTMLを記述する
- [ ] T010 [P] [US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` 内に、月ごとのグループをループ処理し、年月ヘッダーを表示するテンプレートロジックを実装する
- [ ] T011 [P] [US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` 内に、各グループの写真サムネイルをグリッドアイテムとして表示するテンプレートロジックを実装する
- [ ] T012 [US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` に、CSS GridまたはFlexboxを使用して、隙間のないレスポンシブなグリッドレイアウトのスタイルを記述する (T009, T011に依存)

**Checkpoint**: この時点で、User Story 1と2の両方が機能し、ブラウザで最終的なデザインが確認できるはずです。

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: 最終的な品質向上と動作確認

- [ ] T013 コード全体をレビューし、不要なコードの削除やコメントの追加を行う
- [ ] T014 `quickstart.md` に記載された手動テストを実行し、すべての要件が満たされていることを確認する

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 最初に実行する必要があります。
- **User Story 1 (Phase 2)**: Setup完了後に開始できます。
- **User Story 2 (Phase 3)**: User Story 1完了後に開始できます。バックエンドのデータ構造に依存するためです。
- **Polish (Phase 4)**: すべてのユーザーストーリー完了後に実行します。

### User Story Dependencies

- **User Story 2** は **User Story 1** に依存します。

### Within Each User Story

- テストは実装より先に作成し、失敗することを確認します。
- バックエンドロジック (US1) はフロントエンドテンプレート (US2) より先に実装します。

### Parallel Opportunities

- User Story 1のテストタスク (T002, T003, T004) は並行して作成可能です。
- User Story 2のテンプレート実装タスク (T010, T011) は並行して進めることができます。

---

## Implementation Strategy

### MVP First (User Story 1 & 2)

1. Phase 1: Setup を完了します。
2. Phase 2: User Story 1 を完了します (バックエンドロジック)。
3. Phase 3: User Story 2 を完了します (フロントエンド表示)。
4. **STOP and VALIDATE**: 手動テストを実行し、MVPが完全に機能することを確認します。
5. Phase 4: Polish を実行し、最終的な品質を確保します。

この機能はバックエンドとフロントエンドが密接に関連しているため、両方のユーザーストーリーが完了して初めて価値を提供します。したがって、MVPスコープはUS1とUS2の両方を含みます。