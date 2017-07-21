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
	r.LoadHTMLGlob("/home/luren5/Project/src/github.com/luren5/mcat/IDE/templ/*")
	// index
	r.GET("/", index)
	// upload file
	r.POST("/upload-file", uploadFile)
	// edit file
	r.Any("/edit/:fileName", edit)
	// new file
	r.GET("/new-file/:fileName", newFile)

	port, err := utils.Config("ide_port")
	if err != nil {
		r.Run()
	}
	r.Run(":" + port.(string))
	fmt.Println("IDE is on,listening " + port.(string))
}

func init() {
	RootCmd.AddCommand(IDECmd)
}

// index
func index(c *gin.Context) {
	// lis files
	c.HTML(http.StatusOK, "index.templ", gin.H{
		"fileSet": getFileSet(),
	})
}

// edit
func edit(c *gin.Context) {
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
		"fileSet":     getFileSet(),
	})
}

// upload
func uploadFile(c *gin.Context) {
	_, file, err := c.Request.FormFile("new_sol")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": FAIL,
			"msg":    err.Error(),
		})
	}
	if !strings.HasSuffix(file.Filename, ".sol") {
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
}

// new file
func newFile(c *gin.Context) {
	fileName := c.Param("fileName")
	if !strings.HasSuffix(fileName, ".sol") {
		c.JSON(http.StatusOK, gin.H{
			"status": FAIL,
			"msg":    "Invalid file name.",
		})
		return
	}
	newFile := utils.ContractsDir() + fileName
	f, err := os.OpenFile(newFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": FAIL,
			"msg":    err.Error(),
		})
		return
	}
	f.Close()
	c.JSON(http.StatusOK, gin.H{
		"status": SUCCESS,
	})
}

//get file set
func getFileSet() []string {
	files, err := ioutil.ReadDir(utils.ContractsDir())
	if err != nil {
		return []string{}
	}
	var fileSet []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileSet = append(fileSet, f.Name())
	}
	return fileSet
}
