package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

const (
	SUCCESS = iota
	FAIL
)

// IDECmd represents the IDE command
var IDECmd = &cobra.Command{
	Use:   "IDE",
	Short: "Solidity local online IDE.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		startIDE()

		time.Sleep(time.Second * 50)

		fmt.Println("Starting online IDE, listening on 8080…")

	},
}

func startIDE() {
	r := gin.Default()
	r.Static("./static", "./IDE")
	r.LoadHTMLGlob("./IDE/templ/*")
	// index
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.templ", gin.H{
			"title": "hello",
		})
	})

	// upload file
	r.POST("/upload-file", func(c *gin.Context) {
		_, file, err := c.Request.FormFile("new_sol")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": FAIL,
				"msg":    err.Error(),
			})
		}
		if !strings.HasSuffix(file.Filename, ".sol") {
			fmt.Println("file name*******", file.Filename)
			c.JSON(http.StatusOK, gin.H{
				"status": FAIL,
				"msg":    fmt.Sprintf("Invalid file type, %s", file.Filename),
			})
			return
		}

		if f, err := file.Open(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": FAIL,
				"msg":    fmt.Sprintf("Failed to open file, %v", err),
			})
			return
		} else {
			out, err := os.Create(utils.ContractsDir() + "/" + file.Filename)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": FAIL,
					"msg":    fmt.Sprintf("Failed to create file, %v", err),
				})
				return
			}
			defer out.Close()
			if _, err := io.Copy(out, f); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": FAIL,
					"msg":    fmt.Sprintf("Failed to copy file, %v", err),
				})
				return
			}
			// redirect
			c.Redirect(http.StatusTemporaryRedirect, "/edit/"+file.Filename)
		}
	})

	// edit file
	r.Any("/edit/:fileName", func(c *gin.Context) {
		fileName := c.Param("fileName")
		if _, err := os.Stat(utils.ContractsDir() + "/" + fileName); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": FAIL,
				"msg":    fmt.Sprintf("Cant't access to file, %v", err),
			})
			return
		}
		fileContent, err := ioutil.ReadFile(utils.ContractsDir() + "/" + fileName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": FAIL,
				"msg":    fmt.Sprintf("Failed to get file content, %v", err),
			})
			return
		}
		c.HTML(http.StatusOK, "index.templ", gin.H{
			"fileName":    fileName,
			"fileContent": string(fileContent),
		})
	})

	r.Run()
}

func init() {
	RootCmd.AddCommand(IDECmd)
}
