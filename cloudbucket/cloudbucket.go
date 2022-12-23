package cloudbucket

import (
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	storageClient *storage.Client
)

func HandleFileUploadToBucket(c *gin.Context) {
	var err error

	var request FileRequest
	err = c.Bind(&request)

	bucket := c.Request.Form.Get("bucket-name")

	ctx := appengine.NewContext(c.Request)

	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("komerce-be-e1e0765a0e23.json"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	f, uploadedFile, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	defer f.Close()

	sw := storageClient.Bucket(bucket).Object(c.Request.Form.Get("folder-name") + "/" + uploadedFile.Filename).NewWriter(ctx)
	sw.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err := io.Copy(sw, f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	if err := sw.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	u, err := url.Parse("https://storage.googleapis.com/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": u.EscapedPath(),
	})
}

func GetListFile(c *gin.Context) {
	var err error

	var request FileRequest
	err = c.BindJSON(&request)

	bucket := request.BucketName

	ctx := appengine.NewContext(c.Request)

	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("komerce-be-e1e0765a0e23.json"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	var folderName string
	if request.FolderName != "" {
		folderName = request.FolderName + "/"
	}

	query := &storage.Query{Prefix: folderName}
	sw := storageClient.Bucket(bucket)
	var names []string
	it := sw.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, "https://storage.googleapis.com/"+bucket+"/"+attrs.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "request success",
		"files":   names,
	})

}

func GetListFolder(c *gin.Context) {
	var err error

	var request FileRequest
	err = c.BindJSON(&request)

	bucket := request.BucketName

	ctx := appengine.NewContext(c.Request)

	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("komerce-be-e1e0765a0e23.json"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	query := &storage.Query{Prefix: ""}
	sw := storageClient.Bucket(bucket)
	var names []string
	var folderNames []string
	it := sw.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
	}

	for _, s := range names {
		folderName := strings.Split(s, "/")
		if !isAvailable(folderNames, folderName[0]) {
			folderNames = append(folderNames, folderName[0])
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "request success",
		"folder":  folderNames,
	})
}

type FileRequest struct {
	BucketName string `json:"bucket-name" binding:"required"`
	FolderName string `json:"folder-name"`
}

func isAvailable(alpha []string, str string) bool {

	// iterate using the for loop
	for i := 0; i < len(alpha); i++ {
		// check
		if alpha[i] == str {
			// return true
			return true
		}
	}
	return false
}
