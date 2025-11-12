# Quickstart

## spec-kit セットアップ

この機能は、プロジェクトにspec-kitをセットアップし、すぐに利用可能な状態にするためのものです。

### 前提条件

- `git`がインストールされていること。
- `bash`が利用可能な環境であること。

### セットアップ手順

1.  **リポジトリのクローン**:
    ```bash
    git clone https://github.com/tsubasaogawa/lambda-image-viewer.git
    cd lambda-image-viewer
    ```

2.  **セットアップコマンドの実行**:
    `/speckit.specify` コマンドを実行して、spec-kitの初期設定を行います。

    ```bash
    /speckit.specify "spec-kit を運用するために本プロジェクトをスキャンし、spec-kit を稼働させるために必要なファイルを生成して。"
    ```

3.  **生成ファイルの確認**:
    コマンド実行後、プロジェクトのルートディレクトリに`.gemini`ディレクトリが生成され、その中に設定ファイルが格納されていることを確認します。

    ```
    .gemini/
    ├── config.yaml
    └── commands/
        ├── speckit.analyze.toml
        ├── speckit.checklist.toml
        ├── speckit.clarify.toml
        ├── speckit.constitution.toml
        ├── speckit.implement.toml
        ├── speckit.plan.toml
        ├── speckit.specify.toml
        └── speckit.tasks.toml
    ```

### 次のステップ

- `/speckit.plan` を実行して、機能の実装計画を作成します。
- `/spe_ckit.tasks` を実行して、具体的な実装タスクを生成します。
