# Implementation Tasks: カメラ画像の月別グルーピング

**Feature Branch**: `003-cameraroll-monthly-grouping`
**Created**: 2025-11-14

## Phase 1: Setup

- [X] T001 `camera_roll.html.tmpl` の `div.thumbnailBox` に `data-timestamp` 属性を追加し、画像のタイムスタンプを埋め込む `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`

## Phase 2: User Story 1 - 月別での画像表示

**User Story Goal**: ユーザーとして、カメラロールの画像を月ごとにグループ化して表示してほしい。これにより、時系列で写真を簡単に見つけられるようになる。

**Independent Test**: カメラロールページを開き、画像が月ごとにまとまって表示されていれば、このストーリーは達成されたと判断できる。

### Implementation Tasks

- [X] T002 [US1] ページロード時にDOM操作で月別グルーピングを行うJavaScriptの基本構造を追加する `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`
- [X] T003 [US1] `data-timestamp` 属性からタイムスタンプを読み取り、月ごとに画像をグループ化するロジックを実装する `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`
- [X] T004 [US1] グループ化されたデータに基づき、月ごとの見出しと `<ul id="thumbs-YYYY-MM">` を動的に生成し、DOMに挿入するロジックを実装する `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`
- [X] T005 [US1] 元のリストから画像要素を、新しく生成した月別の `<ul>` タグ内に移動させるロジックを実装する `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`
- [X] T006 [US1] 元の `<ul>` タグを空にするか、または削除する処理を追加する `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`
- [X] T007 [US1] グループを降順（新しい月が先頭）にソートして表示するロジックを実装する `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`

## Phase 3: Polish & Cross-Cutting Concerns

- [X] T008 手動テストを実施し、画面のちらつきやJavaScriptエラーがないことを確認する
- [X] T009 [P] タイムスタンプが不正な形式の場合や、`data-timestamp` 属性が存在しない場合のエラーハンドリングを追加する `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`

## Dependencies

- **User Story 1** は **Phase 1** に依存します。

## Parallel Execution

- **Phase 2 (User Story 1)** 内のタスクは、`camera_roll.html.tmpl` という単一のファイルを変更するため、基本的にはシーケンシャルに実行する必要があります。ただし、T003, T004, T005 は密接に関連しているため、単一の関数内でまとめて実装することも可能です。
- **Phase 3** の T009 は独立して実行可能です。

## Implementation Strategy

1.  **MVP (Minimum Viable Product)**: Phase 1 と Phase 2 (User Story 1) のタスクを完了させることで、主要な機能要件を満たすMVPが完成します。
2.  **Incremental Delivery**: まずは基本的なグルーピング機能（T002-T006）を実装し、その後ソート（T007）やエラーハンドリング（T009）を追加していくことで、段階的な開発が可能です。
