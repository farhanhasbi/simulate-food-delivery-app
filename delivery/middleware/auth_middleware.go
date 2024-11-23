package middleware

import (
	"fmt"
	"food-delivery-apps/shared"
	"food-delivery-apps/shared/service"
	"food-delivery-apps/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddlewareWithRole(jwtService service.JwtService, userUc usecase.UserUseCase, allowedRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ensure the JwtService is not nil
		if jwtService == nil{
			shared.SendErrorResponse(ctx, http.StatusUnauthorized, "JWT service not initialized")
			ctx.Abort()
			return
		}

		// Retrieve the Authorization header and ensure it's not empty
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == ""{
			shared.SendErrorResponse(ctx, http.StatusUnauthorized, "Authorization header is missing")
			ctx.Abort()
			return
		}

		// Ensure the Authorization header has the "Bearer " prefix
		if !strings.HasPrefix(authHeader, "Bearer "){
			shared.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid authorization header format")
			ctx.Abort()
			return
		}

		// Extract the token by removing the "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			shared.SendErrorResponse(ctx, http.StatusUnauthorized, "Token is missing after Bearer prefix")
			ctx.Abort()
			return
		}

		// Ensure the token is not blacklist by call IsTokenBlacklisted method
		if userUc.IsTokenBlacklisted(tokenString){
			shared.SendErrorResponse(ctx, http.StatusUnauthorized, "Already logged out")
			ctx.Abort()
			return
		}

		// Validate Token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil{
			shared.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid token: " + err.Error())
			ctx.Abort()
			return
		}

		// Claim Id
		userID, ok := claims["id"].(string)
		if !ok || userID == ""{
			shared.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid or missing user ID in token")
			ctx.Abort()
			return
		}

		// Claim Role
		userRole, ok := claims["role"].(string)
		if !ok || userRole == ""{
			shared.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid or missing user Role in token")
			ctx.Abort()
			return
		}

		// Ensure userRole matches one of the allowedRoles values
		roleAllowed := false
		for _, role := range allowedRoles{
			if strings.EqualFold(userRole, role) {
				roleAllowed = true
				break
			}
		}

		// handle if it's not match
		if !roleAllowed {
			shared.SendErrorResponse(ctx, http.StatusForbidden, fmt.Sprintf("only %s can access this resource", allowedRoles[0]))
			ctx.Abort()
			return
		}

		// Set the user's id, user's role and token in the context for further use in handlers
		ctx.Set("userID", userID)
		ctx.Set("userRole", userRole)
		ctx.Set("tokenString", tokenString)

		ctx.Next()
	}
}