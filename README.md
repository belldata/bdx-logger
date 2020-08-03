# bdx-logger

bdx-logger

## Feature

ログレベル制御が可能なロギング

## How To Use

- ログレベル指定方法

`SetLevel(BxLogLevel)`

| level |
| :--: |
| Debug |
| Info |
| Warning |
| Error | 
| Fatal |
| All |
| Off |

- フォーマット指定方法

`SetLogFormat(string)`

デフォルト: `$date$ $time$ $level$ $prefix$ $file$:$func$:$linenumber$: $message$`

| フォーマット指定文字 | 概要 |
| :-- | :-- |
| `$date$` | 日付(yyyy/mm/dd) |
| `$time$` | 時間(hh:mi:ss) |
| `$level$` | レベル([level]) |
| `$prefix$` | 指定プレフィックス(未指定可) |
| `$file$` | 呼び出しファイル名(ファイル名) |
| `$func$` | 呼び出し関数名 |
| `$linenumber$` | 行番号 |
| `$message$` | 出力内容 |

```go

package main

import "github.com/belldata/bxlogger"

func main() {
    bxlog := bxlogger.New("prefix", bxlogger.Info)
    bxlog.Info("output log message")
}

```
