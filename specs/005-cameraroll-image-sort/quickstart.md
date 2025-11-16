# Quickstart: 全画像を対象とした日付ソート機能

## 概要

この機能は、カメラロールの表示機能を改修し、すべての画像を日付順（降順）でシームレスに閲覧できるようにするものです。バックエンドAPIはGoで実装され、AWS Lambda上で動作します。データストアにはAmazon DynamoDBを使用します。

## 主な変更点

- **`viewer` Lambda関数**:
  - `cameraroll.go` の `CameraRollHandler` を修正し、ページネーションに対応します。
  - クエリパラメータ `last_evaluated_key` と `limit` を受け取り、DynamoDBからのデータ取得時に使用します。
  - `model.ListThumbnails` 関数を呼び出す際の引数を調整します。
- **`model`パッケージ**:
  - `thumbnail.go` の `ListThumbnails` 関数を修正し、ソート順を降順に変更します。
- **フロントエンド**:
  - 無限スクロールを実装し、ユーザーがページ下部に到達したら、次の画像セットをAPIから取得するようにします。
  - APIリクエスト時に `last_evaluated_key` を付与します。

## 開発環境のセットアップ

1.  Go 1.21 と Node.js 18.x をインストールします。
2.  `src/viewer` ディレクトリで `npm install` を実行し、Serverless Frameworkの依存関係をインストールします。
3.  AWS認証情報を設定します。
4.  `serverless deploy` コマンドで、変更をAWS環境にデプロイします。
