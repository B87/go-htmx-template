package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/B87/go-htmx-template/internal/blog"
)

type BlogData struct {
	Posts           []blog.Post
	IsAdmin         bool
	IsAuthenticated bool
}

func BlogHandler(blogRepo blog.BlogRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			posts, err := blogRepo.GetAllPublished()
			if err != nil {
				fmt.Println("Error getting posts:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Convert markdown content to HTML
			for _, post := range posts {
				post.Content = string(MarkdownToHTML(post.Content))
			}
			claims, err := castClaimsFromContext(c)
			if err != nil || claims == nil {
				c.HTML(200, "blog.html", BlogData{Posts: posts, IsAdmin: false, IsAuthenticated: false})
				return
			}
			c.HTML(200, "blog.html", BlogData{Posts: posts, IsAdmin: claims.IsAdmin(), IsAuthenticated: true})
			return
		}
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}

type BlogPostData struct {
	Post            blog.Post
	PostContent     template.HTML
	IsAdmin         bool
	IsAuthenticated bool
}

func BlogPostHandler(blogRepo blog.BlogRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		post, err := blogRepo.GetBySlug(slug)
		if err != nil || post.Status != "published" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		// Convert markdown content to HTML
		postContent := MarkdownToHTML(post.Content)

		claims, err := castClaimsFromContext(c)
		if err != nil || claims == nil {
			c.HTML(200, "blog_post.html", BlogPostData{Post: post, PostContent: postContent, IsAdmin: false, IsAuthenticated: false})
			return
		}
		c.HTML(200, "blog_post.html", BlogPostData{Post: post, PostContent: postContent, IsAdmin: claims.IsAdmin(), IsAuthenticated: true})
	}
}

func BlogAdminHandler(blogRepo blog.BlogRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		claims, err := castClaimsFromContext(c)
		if err != nil || claims == nil || !claims.IsAdmin() {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
		posts, err := blogRepo.GetAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(200, "blog_admin.html", BlogData{Posts: posts, IsAdmin: true, IsAuthenticated: true})
	}
}

func BlogPostAdminHandler(blogRepo blog.BlogRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		claims, err := castClaimsFromContext(c)
		if err != nil || claims == nil || !claims.IsAdmin() {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
		slug := c.Param("slug")

		if c.Request.Method == http.MethodGet {
			post, err := blogRepo.GetBySlug(slug)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.HTML(200, "blog_post_admin.html", BlogPostData{Post: post, IsAdmin: true, IsAuthenticated: true})
			return
		}

		if c.Request.Method == http.MethodPut {
			post, err := blogRepo.GetBySlug(slug)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			post.Title = c.PostForm("title")
			post.Content = c.PostForm("content")
			post.Slug = createSlug(post.Title)
			post.Description = c.PostForm("description")
			post.ThumbnailURL = c.PostForm("thumbnail_url")
			post.Status = c.PostForm("status")
			post.UpdatedAt = time.Now()
			err = blogRepo.Update(post)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			return
		}

		if c.Request.Method == http.MethodPost {
			title := c.PostForm("title")
			content := c.PostForm("content")
			slug := createSlug(title)
			post := blog.Post{
				ID:      uuid.New(),
				Title:   title,
				Slug:    slug,
				Content: content,
				Status:  "draft",
			}
			err := blogRepo.Create(post)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Request.Method = http.MethodGet
			c.Redirect(http.StatusSeeOther, "/admin/blog/posts/"+slug)
			return

		}
	}
}

func createSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.Replace(slug, " ", "-", -1)

	// Remove special characters
	reg, err := regexp.Compile("[^a-zA-Z0-9-]+")
	if err != nil {
		fmt.Println("Regex error:", err)
	}
	slug = reg.ReplaceAllString(slug, "")

	return slug
}
