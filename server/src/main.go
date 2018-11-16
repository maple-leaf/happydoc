package main

import (
	"encoding/base64"
	"crypto/sha256"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions"
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
	sessionKey := os.Getenv("SESSION_KEY")
	store := cookie.NewStore([]byte(sessionKey))
	r.Use(sessions.Sessions("happydoc-session", store))
	r.HTMLRender = services.LoadTemplates("./templates")

	os.Mkdir("documents", 0700)
	r.Static("/documents", "./documents")
	r.Static("/assets", "./assets")

	apiRoutes := r.Group("/document")
	apiRoutes.Use(middlewares.JWTAuth(db))
	{
		apiRoutes.POST("/publish", func(c *gin.Context) {
			status := c.Writer.Status()
			if status == 403 {
				return
			}
			project := c.PostForm("project")
			docType := c.PostForm("type")
			version := c.PostForm("version")
			file, _ := c.FormFile("file")

			zipFilePath := "uploaded/" + file.Filename
			err := c.SaveUploadedFile(file, zipFilePath)

			if err != nil {
				c.Status(500)
				return
			}

			doc := models.Document{
				Project: project,
				CType: docType,
				CVersion: version,
			}

			if (createDocument(doc, zipFilePath, db) != nil) {
				c.Status(500)
				return
			}

			c.Status(200)
		})
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(middlewares.PasswdAuth(db))
	{
		adminGroup.GET("", func(c *gin.Context) {
			accounts := []models.User{}
			db.Find(&accounts)
			c.HTML(200, "admin.html", gin.H{
				"accounts": accounts,
			})
		})

		adminGroup.GET("/new-account", func(c *gin.Context) {
			c.HTML(200, "new-account.html", gin.H{})
		})

		adminGroup.POST("/new-account", func(c *gin.Context) {
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
				c.HTML(200, "account.html", gin.H{
					"name":  t.Name,
					"token": token,
				})
			}
		})

		adminGroup.POST("/login", func(c *gin.Context) {
			if (c.Writer.Status() == 403) {
				c.Done()
				return
			}
			c.Redirect(301, "/admin")
		})

	}


	r.GET("/", func(c *gin.Context) {
		documents := []models.Document{}
		db.Preload("Types").Find(&documents)
		c.HTML(200, "index.html", gin.H{
			"projects": documents,
		})
	})

	r.GET("/setup", func(c *gin.Context) {
		users := []models.User{}
		db.Limit(1).Find(&users)
		if len(users) > 0 {
			c.Redirect(301, "/login")
		}
		c.HTML(200, "setup.html", gin.H{})
	})
	r.POST("/setup", func(c *gin.Context) {
		users := []models.User{}
		db.Limit(1).Find(&users)
		if len(users) > 0 {
			c.Redirect(301, "/login")
		}
		name := c.PostForm("name")
		passwd := c.PostForm("passwd")
		sum := sha256.Sum256([]byte(passwd))
		token := base64.StdEncoding.EncodeToString(sum[:])
		admin := models.User{
			Name: name,
			Token: token,
		}
		x := db.Create(&admin)
		if x.Error != nil {
			c.JSON(500, gin.H{
				"message": "fail to create admin account",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
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
