package utils

//session工具类
import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 初始化系统web session
func InitWebSession() gin.HandlerFunc {

	store := getSessionStore()
	session := sessions.Sessions("ibpdSession", store)
	return session
}

func getSessionStore() cookie.Store {
	sessionMaxAge := 3600
	sessionSecret := "secret"
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge: sessionMaxAge, //seconds
		Path:   "/",
	})
	return store
}
func Session(keyPairs string) gin.HandlerFunc {
	store := getSessionStore()
	return sessions.Sessions(keyPairs, store)
}
