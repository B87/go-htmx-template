package blog

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Post struct {
	ID           uuid.UUID `db:"id"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	Slug         string    `db:"slug"`
	ThumbnailURL string    `db:"thumbnail_url"`
	Content      string    `db:"content"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type PostTranslation struct {
	ID          uuid.UUID `db:"id"`
	PostID      uuid.UUID `db:"post_id"`
	Language    string    `db:"language"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	Content     string    `db:"content"`
	Translation Post      `db:"translation"`
}

type BlogRepository interface {
	GetAll() ([]Post, error)
	GetByID(id uuid.UUID) (Post, error)
	GetBySlug(slug string) (Post, error)
	GetAllPublished() ([]Post, error)
	Create(post Post) error
	Update(post Post) error
	Delete(id uuid.UUID) error
}

var ErrPostNotFound = errors.New("post not found")

type PGBlogRepository struct{ db *sqlx.DB }

func NewPGBlogRepository(db *sqlx.DB) *PGBlogRepository {
	return &PGBlogRepository{db: db}
}

func (r *PGBlogRepository) GetAll() ([]Post, error) {
	var posts []Post
	err := r.db.Select(&posts, "SELECT * FROM blog_posts")
	return posts, err
}

func (r *PGBlogRepository) GetAllPublished() ([]Post, error) {
	var posts []Post
	err := r.db.Select(&posts, "SELECT * FROM blog_posts WHERE status = 'published'")
	return posts, err
}

func (r *PGBlogRepository) GetByID(id uuid.UUID) (Post, error) {
	var post Post
	err := r.db.Get(&post, "SELECT * FROM blog_posts WHERE id = $1", id)
	return post, err
}

func (r *PGBlogRepository) GetBySlug(slug string) (Post, error) {
	var post Post
	err := r.db.Get(&post, "SELECT * FROM blog_posts WHERE slug = $1", slug)
	if err == sql.ErrNoRows {
		return Post{}, ErrPostNotFound
	}
	return post, err
}

func (r *PGBlogRepository) Create(post Post) error {
	_, err := r.db.Exec("INSERT INTO blog_posts (id, title, slug, description, thumbnail_url, content, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		post.ID, post.Title, post.Slug, post.Description, post.ThumbnailURL, post.Content, post.Status, time.Now(), time.Now())
	return err
}

func (r *PGBlogRepository) Update(post Post) error {
	_, err := r.db.Exec("UPDATE blog_posts SET title = $1, slug = $2, description = $3, thumbnail_url = $4, content = $5, status = $6, updated_at = $7 WHERE id = $8",
		post.Title, post.Slug, post.Description, post.ThumbnailURL, post.Content, post.Status, time.Now(), post.ID)
	return err
}

func (r *PGBlogRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM blog_posts WHERE id = $1", id)
	return err
}
