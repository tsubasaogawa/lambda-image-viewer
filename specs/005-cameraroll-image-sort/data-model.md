# Data Model: 全画像を対象とした日付ソート機能

## Entity: Thumbnail (既存の構造体を再利用)

### Attributes
- `Id`: 画像の一意な識別子 (string)
- `Timestamp`: 画像のタイムスタンプ (int64, ソートキーとして使用)
- `Private`: プライベート画像フラグ (bool)
- `Width`: サムネイルの幅 (int32)
- `Height`: サムネイルの高さ (int32)

### Relationships
- なし

### Validation Rules
- `Id` は必須であり、一意であること。
- `Timestamp` は必須であり、有効なUnixタイムスタンプであること。

### State Transitions
- なし
