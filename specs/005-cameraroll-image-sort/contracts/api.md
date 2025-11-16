# API Contract: Camera Roll

## GET /cameraroll

### Summary
カメラロールの画像リストを日付順に取得します。

### Parameters

| Name | In | Description | Required | Schema |
|---|---|---|---|---|
| `last_evaluated_key` | query | ページネーションのためのキー | No | string |
| `limit` | query | 1ページあたりの最大取得件数 | No | integer |

### Responses

| Status Code | Description | Content |
|---|---|---|
| `200 OK` | 成功 | `application/json` |

#### Response Body

```json
{
  "thumbnails": [
    {
      "id": "string",
      "timestamp": "integer",
      "private": "boolean",
      "width": "integer",
      "height": "integer"
    }
  ],
  "last_evaluated_key": "string"
}
```
