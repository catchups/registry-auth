package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (controller AuthController) Test(c *gin.Context) {

	fmt.Println("111111111111")

}
