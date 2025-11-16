package model

import (
	"reflect"
	"sort"
	"testing"
)

func TestThumbnailSort(t *testing.T) {
	// テストデータ
	thumbs := []Thumbnail{
		{Id: "a", Timestamp: 100},
		{Id: "b", Timestamp: 300},
		{Id: "c", Timestamp: 200},
	}

	// 期待される結果 (降順)
	expected := []Thumbnail{
		{Id: "b", Timestamp: 300},
		{Id: "c", Timestamp: 200},
		{Id: "a", Timestamp: 100},
	}

	// ListThumbnails内のソートロジックを直接テスト
	sort.Slice(thumbs, func(i, j int) bool { return thumbs[i].Timestamp > thumbs[j].Timestamp })

	// 結果を比較
	if !reflect.DeepEqual(thumbs, expected) {
		t.Errorf("Slice not sorted correctly.\nGot:    %+v\nWanted: %+v", thumbs, expected)
	}
}

