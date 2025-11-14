# Research: カメラ画像の月別グルーピング

**Feature Branch**: `003-cameraroll-monthly-grouping`
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

## 2. Research Tasks & Findings

本件は既存機能の改修であり、技術的な不確定要素は少ないため、追加の調査は不要と判断します。

- **決定事項**:
  - **実装言語**: JavaScript (クライアントサイド)
  - **テンプレートエンジン**: Go standard library (html/template) はデータ埋め込みのみに使用し、グルーピングロジックはJavaScriptで実装します。
  - **データ構造**: HTMLのカスタムデータ属性 (`data-timestamp`) を利用してタイムスタンプ情報をJavaScriptに渡し、JavaScript内で月別グループを動的に構築します。
