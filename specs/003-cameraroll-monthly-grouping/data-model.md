# Data Model: カメラ画像の月別グルーピング

**Feature Branch**: `003-cameraroll-monthly-grouping`
**Created**: 2025-11-14

## 1. Backend Data Model (No Change)

本機能はフロントエンドでの実装となるため、バックエンドのデータモデル (`src/viewer/internal/model/thumbnail.go` で定義される `model.Thumbnail` 構造体など) に変更はありません。

既存の `model.Thumbnail` 構造体は、引き続き以下の情報を含みます。

- `Id`: 画像の一意なID
- `Timestamp`: 画像のタイムスタンプ
- `Private`: プライベート画像かどうかを示すフラグ
- `Width`: 画像の幅
- `Height`: 画像の高さ

## 2. Frontend Data Handling

フロントエンドのJavaScriptは、HTMLに埋め込まれたカスタムデータ属性 (`data-timestamp`) から画像のタイムスタンプ情報を取得します。この情報を基に、JavaScriptのメモリ上で以下のような論理的なデータ構造を構築し、DOMを操作して表示を更新します。

### MonthlyGroup (JavaScript内での論理的な構造)

- **Key**: `string` (例: "2025-10") - 年と月を表す文字列
- **Value**: `Array<HTMLElement>` - その月に属する画像要素の配列

この論理的な構造は、JavaScriptの `Map` オブジェクトやプレーンなJavaScriptオブジェクトとして実装されます。
