package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	vm "github.com/vmmgr/node/pkg/api/core/vm/v0"
	"log"
	"net/http"
	"strconv"
)

func NodeAPI() {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("failed to connect to qemu")
	}
	defer conn.Close()

	vmh := vm.NewMainHandler(conn)

	router := gin.Default()
	router.Use(cors)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{

			// VM
			v1.POST("/vm", vmh.Add)
			v1.GET("/vm", vmh.GetAll)

			v1.DELETE("/vm/:id", vmh.Delete)
			v1.PUT("/vm/:id", vmh.Update)
			v1.GET("/vm/:id", vmh.Get)

			// VM Status
			v1.PUT("/vm/:id/power", vmh.Startup)
			v1.GET("/vm/:id/power", vmh.GetStatus)
			v1.DELETE("/vm/:id/power", vmh.Shutdown)

			// VM Reset
			v1.PUT("/vm/:id/reset", vmh.Reset)

		}
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(8080), router))
}

func cors(c *gin.Context) {

	//c.Header("Access-Control-Allow-Headers", "Accept, Content-ID, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-ID", "application/json")
	c.Header("Access-Control-Allow-Credentials", "true")
	//c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
