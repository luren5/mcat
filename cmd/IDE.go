package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/compiler"
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
	r.LoadHTMLGlob("IDE/templ/*")
	// index
	r.GET("/", index)
	// upload file
	r.POST("/upload-file", uploadFile)
	// save file
	r.POST("/save-file", saveFile)
	// edit file
	r.Any("/edit/:fileName", edit)
	// new file
	r.GET("/new-file/:fileName", newFile)
	// do compile
	r.POST("/do-compile", doCompile)
	// refresh list
	r.GET("/refresh-list", refreshList)
	// remove file
	r.GET("/remove-file/:fileName", removeFile)

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
	fileSet := getFileSet()
	lastEditFile, _ := utils.LastEditFile()

	for _, v := range fileSet {
		if v == lastEditFile {
			c.Redirect(http.StatusTemporaryRedirect, "/edit/"+lastEditFile)
		}
	}

	// lis files
	c.HTML(http.StatusOK, "index.templ", gin.H{
		"fileSet":     fileSet,
		"projectName": utils.ProjectName(),
	})
}

// edit
func edit(c *gin.Context) {
	fileName := c.Param("fileName")
	var errMes, fileContent string

	if _, err := os.Stat(utils.ContractsDir() + "/" + fileName); err != nil {
		errMes = fmt.Sprintf("Cant't access to file %s, %v", fileName, err)
	}
	if f, err := ioutil.ReadFile(utils.ContractsDir() + "/" + fileName); err != nil {
		errMes = fmt.Sprintf("Failed to get file content, %v", err)
	} else {
		fileContent = string(f)
	}

	if len(errMes) == 0 {
		utils.LastEditFile(fileName)
	}
	c.HTML(http.StatusOK, "index.templ", gin.H{
		"fileName":    fileName,
		"fileContent": strings.Trim(fileContent, " "),
		"fileSet":     getFileSet(),
		"projectName": utils.ProjectName(),
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

//save file
func saveFile(c *gin.Context) {
	fileName := c.PostForm("fileName")
	fileContent := c.PostForm("fileContent")
	err := writeContent(fileName, fileContent)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": SUCCESS,
			"msg":    "Succeed in saving " + fileName,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": FAIL,
			"msg":    err.Error(),
		})
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

// do compile
func doCompile(c *gin.Context) {
	fileName := c.PostForm("fileName")
	fileContent := c.PostForm("fileContent")
	writeContent(fileName, fileContent)

	// wirte content
	contracts, err := compiler.CompileSolidity(solc, utils.ContractsDir()+fileName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": FAIL,
			"msg":    err.Error(),
		})
		return
	}
	var result []map[string]string
	for name, contract := range contracts {
		nameParts := strings.Split(name, ":")
		contractName := nameParts[len(nameParts)-1]
		abi, _ := json.Marshal(contract.Info.AbiDefinition)
		bin := contract.Code
		r := make(map[string]string)
		r["name"] = contractName + ".sol"
		r["bin"] = bin
		r["abi"] = string(abi)
		result = append(result, r)

		// write file
		ioutil.WriteFile(utils.CompiledDir()+"/"+contractName+".abi", abi, 0660)
		ioutil.WriteFile(utils.CompiledDir()+"/"+contractName+".bin", []byte(bin), 0660)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": SUCCESS,
		"msg":    result,
	})
}

// refresh list
func refreshList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": SUCCESS,
		"msg":    getFileSet(),
	})
}

func removeFile(c *gin.Context) {
	fileName := c.Param("fileName")
	os.Remove(utils.ContractsDir() + fileName)
	c.JSON(http.StatusOK, gin.H{
		"status": SUCCESS,
	})
}

// writeContent
func writeContent(fileName, fileContent string) error {
	if _, err := os.Stat(utils.ContractsDir()); err != nil {
		os.MkdirAll(utils.CompiledDir(), 0777)
	}

	return ioutil.WriteFile(utils.ContractsDir()+fileName, []byte(fileContent), 0777)
}

//get file set
func getFileSet() []string {
	files, err := ioutil.ReadDir(utils.ContractsDir())
	if err != nil {
		return []string{}
	}
	var fileSet []string
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sol") {
			continue
		}
		fileSet = append(fileSet, f.Name())
	}
	return fileSet
}
