package main

import (
	"log"

	"github.com/appwrite/sdk-for-go/appwrite"
	_ "github.com/appwrite/sdk-for-go/id"
)

func main() {
  client := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject("6770875d003cd7a26e8e"),
		// appwrite.WithKey("standard_6c4a935251ddd3e1efd8801e2661e6f4f055ffcf550d17d2fb2b082b1d9396ad80ca26f5dc243bef700bd4d6003986bc1d8da323fa9bfa8e00dadee79919910700bc46dcaad30cd1921ea270cbf568b16324837e972947725215d91c51ef20142a5e7b2e182496e886b9e953f1df14f86b48a00c7d0f9a9680c7ac921342383a"),
	)

  databases := appwrite.NewDatabases(client)

  doc, err := databases.ListDocuments(
    "67709112001c053f6cdf",
    "677091380032b5cf769d",
  )

  if err != nil {
    panic(err)
  }

  // println(doc)
  log.Println(doc)

	// config := &atreugo.Config{
	// 	Host: "localhost",
	// 	Port: 8000,
	// }
	// server := atreugo.New(config)
	//
	// server.Path("GET", "/", func(ctx *atreugo.RequestCtx) error {
	// 	println("Hello World")
	// 	return ctx.HTTPResponse("<h1>Atreugo Micro-Framework</h1>")
	// })
	//
	// err := server.ListenAndServe()
	// if err != nil {
	// 	panic(err)
	// }
}
