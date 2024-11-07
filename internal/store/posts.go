package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostsWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
	SELECT id, content, title, user_id, tags, created_at, updated_at, version
	FROM posts
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return &post, nil
}

func (s *PostsStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET content = $1, title = $2, version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostsWithMetadata, error) {
	shouldFilterTags := len(fq.Tags) > 0

	query := `
	SELECT p.id, p.content, p.title, p.user_id, p.tags, p.created_at, p.updated_at, p.version, count(c.id) as comments_count,
	       u.username
	FROM posts p
	LEFT JOIN comments c ON p.id = c.post_id
	LEFT JOIN users u ON p.user_id = u.id
	JOIN followers f ON p.user_id = f.follower_id
	OR p.user_id = $1
	WHERE f.user_id = $1
	AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
	`

	if shouldFilterTags {
		query += "AND p.tags && $5::varchar[]"
	}

	query += `
	GROUP BY p.id, u.username
	ORDER BY p.created_at ` + fq.Sort + `
	LIMIT $2
	OFFSET $3
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var rows *sql.Rows
	var err error

	if shouldFilterTags {
		rows, err = s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	} else {
		rows, err = s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostsWithMetadata
	for rows.Next() {
		var post PostsWithMetadata
		err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.Title,
			&post.UserID,
			pq.Array(&post.Tags),
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Version,
			&post.CommentsCount,
			&post.User.Username,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, post)
	}

	return feed, nil
}
