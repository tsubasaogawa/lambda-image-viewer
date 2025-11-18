# Implementation Plan: カメラロールのページネーション修正

**Feature Branch**: `006-fix-cameraroll-pagination`
**Feature Spec**: [spec.md](./spec.md)

## 1. Technical Context

- **Affected Components**:
  - `src/viewer/cmd/viewer/cameraroll.go`: DynamoDBから画像メタデータを取得し、テンプレートに渡すハンドラロジック。
  - `src/viewer/cmd/viewer/templates/camera_roll.html.tmpl`: 画像を表示し、次ページへのリンクを生成するHTMLテンプレート。
- **Technology Stack**:
  - Go 1.21
  - AWS SDK for Go (DynamoDB)
  - Serverless Framework
- **Dependencies**:
  - AWS DynamoDB: 画像メタデータの保存先。
- **Data Storage**:
  - DynamoDBのテーブルに画像メタデータが格納されている。主キーとソートキーによるクエリが想定される。
- **Integration Points**:
  - なし。Lambda関数内で完結。
- **Security Considerations**:
  - 既存のIAMロールの権限（DynamoDBへのQueryアクセス）で十分。
- **Open Questions / Needs Clarification**:
  - なし。技術的な解決策は明確。

## 2. Constitution Check

- **Principle Compliance**:
  - この修正は既存のコンポーネントのバグ修正であり、プロジェクトの原則に準拠している。
  - テストファーストの原則に従い、まずページネーションの失敗を再現するテストを作成し、次に修正を行う。
- **Gate Evaluation**:
  - **Spec Quality**: [PASS] - 仕様は明確化され、実装に進む準備ができている。
  - **Clarity**: [PASS] - 質問を通じて不明瞭な点が解消された。

## 3. Phase 0: Outline & Research

- **Objective**: ページネーション問題の最適な技術的解決策を決定する。
- **Key Decisions**:
  - DynamoDBの `Query` APIが返す `LastEvaluatedKey` を利用した、カーソルベースのページネーションを実装する。これは、大規模なデータセットに対して効率的で一貫性のあるページ分割を提供するための標準的な方法である。
- **Artifacts**:
  - [research.md](./research.md)

## 4. Phase 1: Design & Contracts

- **Objective**: データ構造、APIの規約、および開発の開始手順を定義する。
- **Key Designs**:
  - **Data Model**: `Image`エンティティの構造を定義する。ソート順序を安定させるため、プライマリキーを`YYYY-MM`（パーティションキー）、セカンダリキーを`timestamp#id`（ソートキー）のような複合キーとすることが考えられる。
  - **API Contracts**: ページネーションの状態を渡すため、URLクエリパラメータ（例: `?next_cursor=...`）を定義する。`next_cursor`の値は、`LastEvaluatedKey`をBase64エンコードした文字列とする。
- **Artifacts**:
  - [data-model.md](./data-model.md)
  - [contracts/api.md](./contracts/api.md)
  - [quickstart.md](./quickstart.md)

## 5. Phase 2: Implementation & Testing (TBD)

- **Objective**: コードを実装し、テストを通じて品質を保証する。
- **Workflow**:
  1. `cameraroll_test.go`に、現在のページネーションの問題を再現するテストケースを追加する。
  2. `cameraroll.go`を修正し、`LastEvaluatedKey`を使ったカーソルベースのページネーションを実装する。
  3. `camera_roll.html.tmpl`を修正し、次ページへのリンクにカーソルを埋め込む。
  4. すべてのテストがパスすることを確認する。
- **Artifacts**:
  - `tasks.md` (TBD)