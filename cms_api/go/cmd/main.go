package main

import (
	"cms_api/internal/config"
	route "cms_api/internal/di"
	"log"
)

func main() {
	// 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("スタンドアロンサーバー設定の読み込みに失敗しました: %v", err)
	}

	log.Printf("スタンドアロンサーバーを初期化します: DB=%s:%d/%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// Echo サーバーの初期化と起動
	e := route.RouteHandler(cfg)

	// サーバーアドレスの設定
	address := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("CMS APIサーバーを開始します: http://%s", address)

	// サーバー開始（ブロッキング）
	if err := e.Start(address); err != nil {
		log.Fatalf("サーバーの開始に失敗しました: %v", err)
	}
}
