# Tasks: Camera Roll デザイン改修 (フロントエンドのみ)

**Input**: Design documents from `/specs/004-cameraroll-design-update/`
**Prerequisites**: plan.md, spec.md

**方針**: バックエンドのGoコードは変更せず、フロントエンドのHTML, CSS, JavaScriptのみでデザイン改修を完結させます。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 並行して実行可能
- **[Story]**: 関連するユーザーストーリー
- 説明には正確なファイルパスを含める

---

## Phase 1: HTML構造の刷新 (US1, US2)

**Goal**: 新しいグリッドレイアウトの基礎となるHTML構造を準備する。

- [x] T001 [US1, US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の既存の `<ul>` タグを、グリッドコンテナ用の新しい `<div>` (例: `<div id="cameraroll-grid"></div>`) に置き換える。
- [x] T002 [US1, US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` 内の既存の月別グループ化JavaScriptを削除する。

---

## Phase 2: CSSによるグリッドレイアウトの実装 (US2)

**Goal**: FlickrのようなレスポンシブなグリッドレイアウトをCSSで実現する。

- [ ] T003 [P] [US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の `<style>` タグ内に、グリッドコンテナ用のCSSを追加する。(`display: grid`, `grid-template-columns`, `gap`など)
- [ ] T004 [P] [US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の `<style>` タグ内に、メディアクエリを使用して、異なる画面幅に応じたカラム数を定義するレスポンシブCSSを追加する。
- [ ] T005 [P] [US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の `<style>` タグ内に、月ヘッダー (`<h2>`など) のためのスタイルを追加する。

---

## Phase 3: JavaScriptによるロジック実装 (US1, US2)

**Goal**: JavaScriptを使用して、写真の月別グループ化とDOMへの描画を行う。

- [ ] T006 [US1] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の `<script>` タグ内に、ページ読み込み時に `{{ .Thumbnails }}` から全サムネイルの情報をJavaScriptの配列に変換する処理を実装する。
- [ ] T007 [US1] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の `<script>` タグ内に、サムネイル配列をタイムスタンプに基づいて月ごとにグループ化する関数を実装する。タイムスタンプがない場合は「Undefined」グループに分類する。
- [ ] T008 [US1] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の `<script>` タグ内に、グループ化したデータを月（新しい順）でソートする処理を実装する。「Undefined」は最後に表示する。
- [ ] T009 [US2] `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl` の `<script>` タグ内に、ソート済みのグループデータに基づいて、月のヘッダーと写真のサムネイルをグリッドコンテナ内に動的に生成し、DOMに追加する処理を実装する。(T001, T007, T008に依存)

---

## Phase 4: 仕上げと確認

**Purpose**: 最終的な品質向上と動作確認

- [ ] T010 コード全体をレビューし、不要なコードの削除やコメントの追加を行う。
- [ ] T011 `quickstart.md` に記載された手動テストを実行し、すべての要件が満たされていることを確認する。
