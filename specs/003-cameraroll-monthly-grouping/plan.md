# Implementation Plan: カメラ画像の月別グルーピング

**Feature Branch**: `003-cameraroll-monthly-grouping`
**Feature Spec**: [spec.md](./spec.md)
**Created**: 2025-11-14

## 1. Technical Context

本機能は、クライアントサイドのJavaScriptとDOM操作によって実現します。

- **対象ファイル**:
  - `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`: カメラロールのHTML構造を定義するGoテンプレートです。ここにJavaScriptコードを追加し、各サムネイルにタイムスタンプ情報を埋め込みます。

- **既存ロジック**:
  - `cameraroll.go` の `generateCamerarollHtml` 関数は、DynamoDBからサムネイル情報のリスト (`[]model.Thumbnail`) を取得し、`CameraRollData` 構造体に格納して `camera_roll.html.tmpl` テンプレートに渡します。この部分は**変更しません**。
  - 現在のテンプレートは、単一の `<ul>` タグ内で画像のリストをループ処理して表示しています。

- **改修方針**:
  1. `camera_roll.html.tmpl` を修正し、各サムネイルの `div.thumbnailBox` 要素に、JavaScriptで読み取れる形式（例: `data-timestamp` 属性）で画像のタイムスタンプを埋め込みます。
  2. `camera_roll.html.tmpl` にJavaScriptコードを追加します。
  3. 追加するJavaScriptコードは、ページロード後に以下の処理を実行します。
     a. すべての画像要素（またはその親要素）を取得します。
     b. 各要素からタイムスタンプ情報を抽出し、画像を年と月でグループ化します。
     c. 月ごとの見出し（例: `<h2>2025年10月</h2>`）と、`<ul id="thumbs-YYYY-MM">` 形式のIDを持つ `<ul>` タグを動的に生成します。
     d. 生成した見出しと `<ul>` タグをDOMに挿入し、元のリストから画像要素を適切な月の `<ul>` タグの中に移動させます。

## 2. Constitution Check

本プロジェクトの憲法（`constitution.md`）は現在プレースホルダーであり、具体的な原則は定義されていません。したがって、この実装が憲法に違反することはありません。

## 3. Phase 0: Outline & Research

本件は既存機能の改修であり、技術的な不確定要素は少ないですが、フロントエンドでのDOM操作のパフォーマンスとユーザー体験への影響について考慮します。

- **決定事項**:
  - **実装言語**: JavaScript (クライアントサイド)
  - **テンプレートエンジン**: Go standard library (html/template) はデータ埋め込みのみに使用し、グルーピングロジックはJavaScriptで実装します。
  - **データ構造**: HTMLのカスタムデータ属性 (`data-timestamp`) を利用してタイムスタンプ情報をJavaScriptに渡し、JavaScript内で月別グループを動的に構築します。

調査結果の詳細は `research.md` に記述します。

## 4. Phase 1: Design & Contracts

### データモデル (`data-model.md`)

バックエンドの `model.Thumbnail` 構造体は変更しません。フロントエンドのJavaScriptは、HTMLに埋め込まれた `data-timestamp` 属性からタイムスタンプ情報を取得し、これを基に月別グループをメモリ上で構築します。

### APIコントラクト (`/contracts`)

本改修はクライアントサイドの表示ロジックの変更であり、外部に公開するAPIの変更は伴わないため、APIコントラクトの作成は不要です。

### クイックスタート (`quickstart.md`)

ローカル環境で変更を確認するための手順を記述します。

## 5. Phase 2: Implementation & Testing

- **タスク分割**:
  - `camera_roll.html.tmpl` のHTML構造とJavaScriptコードの追加
  - JavaScriptによる月別グルーピングロジックの実装
  - フロントエンドの動作確認（手動テスト）
- **テスト方針**:
  - 主に手動での動作確認を行います。ブラウザの開発者ツールを用いて、生成されたHTML構造が仕様通りであること、およびJavaScriptエラーが発生しないことを確認します。
  - 必要に応じて、JavaScriptの単体テストフレームワーク（例: Jest, Mocha）を導入することも検討しますが、今回は小規模なDOM操作であるため、手動テストを主とします。