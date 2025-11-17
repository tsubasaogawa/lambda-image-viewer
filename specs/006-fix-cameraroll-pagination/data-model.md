# Data Model: カメラロールのページネーション修正

## 1. 概要

このドキュメントは、「カメラロールのページネーション修正」機能で扱われる主要なデータエンティティの構造を定義します。既存の `src/viewer/internal/model/thumbnail.go` に定義されている `Thumbnail` 構造体をそのまま利用します。

## 2. エンティティ定義

### Thumbnail (既存)

カメラロールに表示される各画像のメタデータを表します。

**ソース**: `src/viewer/internal/model/thumbnail.go`

**属性**:

| 属性名 | データ型 | 説明 |
| :--- | :--- | :--- |
| `Id` | String | 画像のユニークID（S3のオブジェクトキー）。 |
| `Timestamp` | int64 | 撮影日時のUnixタイムスタンプ（秒）。 |
| `Private` | bool | 写真がプライベートかどうかを示すフラグ。 |
| `Width` | int32 | 元画像の幅。 |
| `Height` | int32 | 元画像の高さ。 |

**DynamoDBのキー構造（既存のコードからの推測）**:

- **テーブルインデックス**: `Timestamp`
  - `ListThumbnails`関数で `Timestamp` インデックスが利用されており、時間でのソートが行われています。
- **ページネーション**:
  - `ListThumbnails`関数は `LastEvaluatedKey` を返しており、カーソルベースのページネーションが実装されています。

## 3. データアクセスパターンの修正方針

- **`ListThumbnails` 関数の修正**:
  - 現在の実装では、DynamoDBから取得した後にGoのコード内でソート (`sort.Slice`) を行っています。これを、DynamoDBの `Query` または `Scan` の `ScanIndexForward: false` を利用して、データベース側で直接降順ソートするように修正します。これにより、パフォーマンスが向上し、ページネーションの一貫性がより厳密に保証されます。
  - `Id` を第二のソートキーとして利用するロジックを追加し、タイムスタンプが同一の場合の順序を安定させます。これはDynamoDBのソートキーの設計変更を伴う可能性があります（例: `timestamp#id`）。