# api/Dockerfile
FROM golang:1.23 as builder

WORKDIR /app

# ファイルをすべてコピー（go.modとgo.sumを含む）
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./server.go

# 実行用の軽量イメージ
FROM alpine:latest

WORKDIR /app

# 必要なパッケージのインストール
RUN apk --no-cache add ca-certificates

# ビルドしたバイナリを明示的に指定してコピー
COPY --from=builder /app/main /app/main

EXPOSE 5002

# 実行権限を確認してから実行
CMD ["./main"]