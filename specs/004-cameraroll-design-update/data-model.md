# Data Model: Camera Roll デザイン改修

この機能改修において、既存のデータモデル（DynamoDBのテーブルスキーマ）に変更はありません。

## Existing Entities

### 写真 (Photo / Item)

- **Source**: `src/viewer/internal/model/item.go`
- **Description**: ユーザーがアップロードした画像ファイルとそのメタデータを格納します。
- **Key Attributes**:
    - `Id` (String): 写真のユニークID
    - `Timestamp` (Number): 写真の撮影日時（Unixタイムスタンプ）
    - `Path` (String): S3上のオブジェクトパス
    - `Thumbnail` (String): サムネイル画像のS3パス
    - `IsPrivate` (Boolean): 非公開フラグ

### 写真グループ (Photo Group)

- **Description**: このエンティティは、アプリケーションロジック内で動的に生成されるものであり、永続化されたデータモデルではありません。`cameraroll.go` の中で、同じ撮影年月を持つ `Photo` の集まりとして扱われます。