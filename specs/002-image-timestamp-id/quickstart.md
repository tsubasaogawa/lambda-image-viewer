# Quickstart: 画像へのタイムスタンプID付与

## 概要

この機能は、Camera Rollページの各画像に、その画像のタイムスタンプを`id`として持つ`<div>`タグを付与します。

## 変更点

- **ファイル**: `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`
- **変更内容**:
  - `range`ループ内の`<div class="thumbnailBox">`タグに`id="{{.Timestamp}}"`を追加。

## テスト方法

1. アプリケーションをデプロイします。
2. Camera Rollページにアクセスします。
3. ブラウザの開発者ツールで、いずれかの画像の要素を検証します。
4. `<div class="thumbnailBox">`タグに、`id`属性が設定されており、その値がタイムスタンプ形式（例: `2023-10-27T07:23:31Z`）であることを確認します。
