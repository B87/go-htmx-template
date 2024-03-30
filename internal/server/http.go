package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/B87/go-htmx-template/internal/auth"
	"github.com/B87/go-htmx-template/internal/blog"
	"github.com/B87/go-htmx-template/internal/server/handlers"
)

const TEMPLATES_DIR = "web/templates/%s"
const STATIC_DIR = "web/static"

type HttpServer struct {
	engine *gin.Engine
	srv    *http.Server
	CDN    CDN
	UseCDN bool
}

type HTTPServerConfig struct {
	Host   string
	Port   int
	CDN    CDN
	UseCDN bool
	DB     *sqlx.DB
}

func NewHttpServer(config HTTPServerConfig) *HttpServer {

	router := gin.Default()

	router.Use(AuthMiddleware())

	// Load HTML templates
	router.LoadHTMLGlob("web/templates/*.html")

	userRepo := auth.NewPostgresUserRepository(config.DB)
	blogRepo := blog.NewPGBlogRepository(config.DB)

	// Setup routes

	registerStaticRoutes(router, config.CDN, config.UseCDN)

	router.GET("/", handlers.HomeHandler)

	router.GET("/login", handlers.LoginHandler(userRepo))
	router.POST("/login", handlers.LoginHandler(userRepo))

	router.GET("/signup", handlers.SignupHandler(userRepo))
	router.POST("/signup", handlers.SignupHandler(userRepo))

	router.GET("/logout", handlers.LogoutHandler)

	router.GET("/admin", handlers.AdminHandler())
	router.GET("/admin/settings", handlers.AdminSettings())

	router.GET("/blog", handlers.BlogHandler(blogRepo))
	router.GET("/blog/:slug", handlers.BlogPostHandler(blogRepo))

	router.GET("/admin/blog", handlers.BlogAdminHandler(blogRepo))
	router.GET("/admin/blog/posts/:slug", handlers.BlogPostAdminHandler(blogRepo))
	router.PUT("/admin/blog/posts/:slug", handlers.BlogPostAdminHandler(blogRepo))
	router.POST("/admin/blog/posts", handlers.BlogPostAdminHandler(blogRepo))

	return &HttpServer{
		engine: router,
		srv: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
		CDN: config.CDN,
	}
}

func (server *HttpServer) Serve() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("Shutting down gracefully")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

// Register routes for serving static files, either from a CDN or from localhost
func registerStaticRoutes(engine *gin.Engine, cdn CDN, useCDN bool) {
	if useCDN {
		log.Println("Using CDN for static files")
		engine.Static("/static", cdn.RootURL())
	} else {
		log.Println("Using localhost for static files")
		engine.Static("/static", STATIC_DIR)
	}
}

// Upload static files to the CDN
func (server *HttpServer) UploadStaticFiles() {
	if server.CDN == nil {
		log.Println("No CDN configured, skipping static files upload")
		return
	}

	err := server.CDN.UploadFolder(STATIC_DIR)
	if err != nil {
		log.Fatalf("Failed to upload static files: %v", err)
	}
}

// Middleware for parsing JWT cookie
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/signup" || strings.Contains(c.Request.URL.Path, "/static") {
			c.Next()
		}
		// Check if the JWT cookie is present
		cookie, err := c.Cookie("jwt")
		if err == http.ErrNoCookie || cookie == "" {
			// If the cookie is not set, set IsAuthenticated to false and continue
			c.Set("IsAuthenticated", false)
			c.Next()
		}
		// If the JWT cookie is present, validate the token
		claims, err := auth.ValidateJWTToken(cookie)
		switch {
		case errors.Is(err, auth.ErrEmpotyToken):
			c.Set("IsAuthenticated", false)
			c.Next()
		case errors.Is(err, auth.ErrInvalidTokenFormat):
			log.Println("Invalid token format:", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		case errors.Is(err, auth.ErrInvalidSignature):
			log.Println("Invalid token signature:", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		case errors.Is(err, auth.ErrExpiredToken):
			log.Println("Token expired, redirecting to login")
			c.Request.Method = http.MethodGet
			c.Redirect(http.StatusSeeOther, "/login")
		case err != nil && !errors.Is(err, auth.ErrExpiredToken):
			log.Println("Failed to parse JWT token:", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("IsAuthenticated", true)
		c.Set("Claims", claims)
		c.Next()
	}
}
