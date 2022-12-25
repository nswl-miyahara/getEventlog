# PC起動ログ取得ツール

## 概要
- PC（Windows）の起動ログを取得するツール
- Powershellでイベントログを取得してTSVに出力する

## ビルド
```
GOOS=windows GOARCH=amd64 go build -o getEventLog.exe main.go
```

## 使い方
- .envの変更（getEventLog.exeと同じ階層に置く）
    ```
    OUTPUT=<出力ファイル名>
    ```
- Powershellで実行
    ```
    PS C:\Users\YRL-0510\Desktop> .\getEventLog.exe
    ```
