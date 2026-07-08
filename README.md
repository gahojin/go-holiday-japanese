# go-holiday-japanese

日本の祝日判定ユーティリティ for Go

[![GoDoc](https://godoc.org/github.com/gahojin/go-holiday-japanese?status.svg)](https://godoc.org/github.com/gahojin/go-holiday-japanese)
[![Go Report Card](https://goreportcard.com/badge/github.com/gahojin/go-holiday-japanese)](https://goreportcard.com/report/github.com/gahojin/go-holiday-japanese)
[![License](https://img.shields.io/github/license/gahojin/go-holiday-japanese)](LICENSE)

## 使い方

### インストール

```bash
go get github.com/gahojin/go-holiday-japanese
```

### サンプルコード

```go
package main

import (
	"fmt"
	"time"

	"github.com/gahojin/go-holiday-japanese"
)

func main() {
	// 祝日かどうかの判定
	t := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if holiday.IsHoliday(t) {
		fmt.Println("今日は祝日です")
	}

	// 祝日名の取得
	if name := holiday.GetHolidayName(t); name != nil {
		fmt.Printf("祝日名: %s (%s)\n", name.Ja, name.En)
	}

	// 期間内の祝日一覧を取得
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	holidays := holiday.Between(start, end)
	for _, h := range holidays {
		fmt.Printf("%s: %s\n", h.Date.Format("2006-01-02"), h.Name.Ja)
	}
}
```

## ライセンス

[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0)

```
Copyright 2025, GAHOJIN, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
