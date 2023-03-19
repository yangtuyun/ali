package main

import "log"

func main() {
	//                Ir.NONCE_MIN = 0,
	//                Ir.NONCE_MAX = 2147483647;
	//sign()
	//alistSession()
	url := AlistDownload("")
	log.Println(url)
	//app := fiber2.New()
	//app.Use(cors.New())
	//app.Get("/", func(ctx *fiber2.Ctx) error {
	//
	//	err := InitAliKey()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	err = CreateSession()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	url := AliDownload("")
	//	return ctx.Redirect(url)
	//
	//})
	//log.Fatal(app.Listen(":3000"))
}
