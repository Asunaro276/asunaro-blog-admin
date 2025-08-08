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
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// Echo サーバーの初期化と起動
	e := route.RouteHandler(cfg)

	// サーバーアドレスの設定
	address := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("CMS APIサーバーを開始します: http://%s", address)

	e.Logger.Fatal(e.Start(address))
}
