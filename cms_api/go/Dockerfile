FROM golang:1.24-alpine

WORKDIR /app

# 必要なパッケージをインストール
RUN apk add --no-cache git

# 依存関係をコピーしてインストール
COPY admin_api_server/go.mod admin_api_server/go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ローカル開発用サーバーをビルド
RUN cd admin_api_server && go build -o /app/server cmd/local/main.go

EXPOSE 8080

# サーバー起動
CMD ["/app/server"] 
