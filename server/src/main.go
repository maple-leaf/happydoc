package main

import (
	"os"
	"compress/flate"
	"github.com/mholt/archiver"
	"github.com/maple-leaf/happydoc-server/models"

	"github.com/maple-leaf/happydoc-server/services"

	"github.com/maple-leaf/happydoc-server/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type App struct {
	Db *gorm.DB
}

var app = App{}

func main() {
	psqlPassWD := os.Getenv("DB_PASSWD")
	psql := "host=db port=5432 user=postgres dbname=postgres password=" + psqlPassWD + " sslmode=disable"
	db, err := gorm.Open("postgres", psql)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	app.Db = db

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Document{})
	db.AutoMigrate(&models.DocumentType{})

	router := setupRoutes(db)
	router.Run() // listen and serve on 0.0.0.0:8080
}

func setupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.HTMLRender = services.LoadTemplates("./templates")

	os.Mkdir("documents", 0700)
	r.Static("/documents", "./documents")
	r.Static("/assets", "./assets")

	apiRoutes := r.Group("/document")
	apiRoutes.Use(middlewares.JWT(db))
	{
		apiRoutes.POST("/publish", func(c *gin.Context) {
			project := c.PostForm("project")
			docType := c.PostForm("type")
			version := c.PostForm("version")
			file, _ := c.FormFile("file")

			zipFilePath := "uploaded/" + file.Filename
			err := c.SaveUploadedFile(file, zipFilePath)

			if err != nil {
				c.Status(500)
			}

			doc := models.Document{
				Project: project,
				CType: docType,
				CVersion: version,
			}

			if (createDocument(doc, zipFilePath, db) != nil) {
				c.Status(500)
			}

			c.Status(200)
		})
	}

	r.POST("/token/new", func(c *gin.Context) {
		t := models.User{}
		c.ShouldBind(&t)
		token, publicKey, err := services.GenerateJWTRS(services.AuthClaims{
			Username: t.Name,
		})
		if err != nil {
			panic(err)
		}

		t.Token = publicKey
		err = services.DB{Db: db}.Insert(&t)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Failed to generate new token",
			})
		} else {
			c.JSON(200, gin.H{
				"name":  t.Name,
				"token": token,
			})
		}
	})

	r.GET("/", func(c *gin.Context) {
		documents := []models.Document{}
		db.Preload("Types").Find(&documents)
		c.HTML(200, "index.html", gin.H{
			"projects": documents,
		})
	})

	return r
}

func createDocument(doc models.Document, zipFilePath string, db *gorm.DB) (err error) {
		docFolderPath := "documents/" + doc.Project + "/" + doc.CType + "/" + doc.CVersion
		err = os.MkdirAll(docFolderPath, 0700)

		if err != nil {
			return
		}

		z := archiver.Zip{
			CompressionLevel:       flate.DefaultCompression,
			MkdirAll:               true,
			SelectiveCompression:   true,
			ContinueOnError:        false,
			OverwriteExisting:      false,
			ImplicitTopLevelFolder: false,
		}
		err = z.Unarchive(zipFilePath, docFolderPath)

		os.Remove(zipFilePath)
		err = doc.CreateOrUpdate(db)

		return
}
