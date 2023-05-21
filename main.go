package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

type image struct {
	ID int `form:"id"`
	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
}

var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
    router := gin.Default()
    router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.POST("/image", postImage)
	router.POST("/images", uploadSingleFile)
    router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	err := c.BindJSON(&newAlbum)
	if  err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, albums)
}

func postImage(c *gin.Context) {
	// var imageObj image
	// if err := c.ShouldBind(imageObj); err != nil {
	// 	c.String(http.StatusBadRequest, "bad request")
	// 	return
	// }

	// err := c.SaveUploadedFile(imageObj.Avatar, "/assets"+imageObj.Avatar.Filename)

	// if err != nil {
	// 	c.String(http.StatusBadRequest, "bad request")
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"status":"ok",
	// 	"data":imageObj,
	// })
	file, err := c.FormFile("dragonball1")

	if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "message": "No file is received",
        })
        return
    }

	extension := filepath.Ext(file.Filename)
    newFileName := "dragonball1" + extension

	out, err := os.Create("Z:\\projects\\go-practice\\images\\"+newFileName)

	if err != nil {
		fmt.Println("error")
	}

	defer out.Close()

    // The file is received, so let's save it
    // if err := c.SaveUploadedFile(file, "/some/path/on/server/" + newFileName); err != nil {
    //     c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
    //         "message": "Unable to save the file",
    //     })
    //     return
    // }

	fmt.Println("success ", newFileName)


    // File saved successfully. Return proper result
    c.JSON(http.StatusOK, gin.H{
        "message": "Your file has been successfully uploaded.",
    })


}	

func uploadSingleFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	filePath := "Z:\\projects\\go-practice\\images\\" + filename

	out, err := os.Create("Z:\\projects\\go-practice\\images\\" + filename)
	if err != nil {
		fmt.Println("error")
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println("error")
	}
	ctx.JSON(http.StatusOK, gin.H{"filepath": filePath})
}