package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Pessoa struct {
	Id    int    `json:"id" gorm:"primaryKey"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

var db = gorm.DB{}

func main() {

	dbURL := "postgres://postgres:postgres@localhost:5432/go"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		fmt.Printf("erro ---------------------")
	}

	db.AutoMigrate(&Pessoa{}) //cria a tabela

	fmt.Println("Iniciando:")

	r := gin.Default()

	r.GET("/pessoa/:usuario", func(c *gin.Context) {
		p := c.Params.ByName("usuario")

		var pessoas []Pessoa
		fmt.Println("Buscando...")
		db.Where("nome = ?", p).Find(&pessoas)

		c.JSON(http.StatusOK, &pessoas)
	})

	r.GET("/pessoa", func(c *gin.Context) {

		var pessoas []Pessoa
		fmt.Println("Buscando...")
		db.Find(&pessoas)

		c.JSON(http.StatusOK, &pessoas)
	})

	r.POST("pessoa", func(c *gin.Context) {

		var pessoa Pessoa

		if c.Bind(&pessoa) == nil {

			encoder := json.NewEncoder(os.Stdout)
			encoder.Encode(pessoa)

			fmt.Println(" - ", pessoa.Nome)

			db.Create(&pessoa)

			c.JSON(http.StatusCreated, gin.H{"status": "ok"})
		} else {
			fmt.Println("erro")
		}
	})

	r.DELETE("/pessoa/:id", func(c *gin.Context) {

		id := c.Params.ByName("id")

		var p Pessoa
		db.Where("id = ?", id).First(&p)

		db.Delete(p)

		c.JSON(http.StatusOK, "Removido")
	})

	r.PUT("/pessoa/:id", func(c *gin.Context) {

		id := c.Params.ByName("id")

		var pessoa Pessoa

		if c.Bind(&pessoa) == nil {

			var old Pessoa
			db.Where("id = ?", id).First(&old)

			encoder := json.NewEncoder(os.Stdout)
			encoder.Encode(pessoa)

			old.Email = pessoa.Email
			old.Nome = pessoa.Nome

			db.Save(&old)

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			fmt.Println("erro")
		}
	})

	r.Run(":8080")

}
