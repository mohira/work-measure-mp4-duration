# 特定ディレクトリ直下のmp4の再生時間を出力するだけ

## 使い方

```
target_dir
├── other_file
└── other_dir
    ├── otherA.mp4
    └── otherB.mp4
├── sample01.mp4
└── sample02.mp4
```

```
$ go run main.go target_dir
01:11	sample01.mp4
02:22	sample02.mp4
合計時間 03:33
```

```
# 直下にmp4がない場合
$ go run main.go ~/Documents/                               
2021/05/14 15:45:51 /Users/mohira/Documents/直下にはmp4ファイルは1つもなかったよ
exit status 1
```