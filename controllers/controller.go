package controllers

import (
	"fmt"
	"io"
	userModel "main_gateway/models/user"
	"main_gateway/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	if ctx.Request.ContentLength == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"data":  "Request body is empty",
		})
		return
	}
	var LoginDetails userModel.LoginDetails

	if err := ctx.BindJSON(&LoginDetails); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "data": err.Error()})
		ctx.Abort()
		return
	}
	// Check if the username and password are correct or not in DB
	userID, err := userModel.AuthenticateUser(LoginDetails.Username, LoginDetails.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "data": "Username Or Password is incorrect"})
		ctx.Abort()
		return
	}
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "data": "Username Or Password is incorrect"})
		ctx.Abort()
		return
	}

	user, err := userModel.GetUser(userID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "data": "User Details Not Found"})
		ctx.Abort()
		return
	}
	JwtKey, err := userModel.GenrateJwtWithClaims(*user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "data": "error genrating Token"})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"error": false, "authenticated": true, "data": JwtKey})
}

func ProxyToService(route utils.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request method is allowed for this route
		if !utils.Contains(route.Methods, c.Request.Method) {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}

		// Construct the target URL including query parameters
		targetURL := route.Target + c.Param("path") + "?" + c.Request.URL.RawQuery
		fmt.Println(targetURL)

		// Create a new request to the target service
		targetReq, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
			return
		}

		// Copy headers from the original request to the target request
		for key, values := range c.Request.Header {
			for _, value := range values {
				targetReq.Header.Add(key, value)
			}
		}

		// Copy query parameters from the original request to the target request
		targetReq.URL.RawQuery = c.Request.URL.RawQuery

		// Copy form parameters from the original request to the target request
		c.Request.ParseForm()
		for key, values := range c.Request.PostForm {
			for _, value := range values {
				targetReq.Form.Add(key, value)
			}
		}

		// Send the request to the target service
		client := &http.Client{}
		resp, err := client.Do(targetReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request to target service"})
			return
		}
		defer resp.Body.Close()

		// Copy the response headers from the target response to the original response
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Copy the response status code and body from the target response to the original response
		c.Status(resp.StatusCode)
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying response body"})
		}
	}
}
