package repository

import (
	"blog/config"
	"blog/internal/model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(conf config.Config) (*Repository, error) {
	url := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=disable",
		conf.Host,
		conf.Username,
		conf.DB,
		conf.Password,
	)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %v", err)
	}

	return &Repository{db: db}, nil
}

func (repo Repository) Migrate(models ...interface{}) error {
	return repo.db.AutoMigrate(models...)
}

func (repo Repository) Create() error {
	var opt int
	fmt.Print(
		"  1. Create post\n",
		"  2. Create comment\n",
		"Choose action: ",
	)

	if _, err := fmt.Scanf("%d\n", &opt); err != nil {
		var temp string
		fmt.Scanln(&temp)
	}

	switch opt {
	case 1:
		if err := repo.createPost(); err != nil {
			return fmt.Errorf("repo.createPost: %v", err)
		}
	case 2:
		if err := repo.createComment(); err != nil {
			return fmt.Errorf("repo.createPost: %v", err)
		}
	default:
		fmt.Printf("Incorrect input\n\n")
	}
	return nil
}

func (repo Repository) createPost() error {
	var (
		title   string
		content string
	)
	fmt.Print("Enter post Title: ")
	fmt.Scanln(&title)
	fmt.Print("Enter post Content: ")
	fmt.Scanln(&content)

	post := &model.Post{
		Title:   title,
		Content: content,
	}
	result := repo.db.Create(post)

	if result.Error != nil {
		return fmt.Errorf("repo.db.Create: %v", result.Error)
	}

	fmt.Printf("Post created successfully!\n\n")
	return nil
}

func (repo Repository) createComment() error {
	var (
		post    model.Post
		content string
	)

	postID, err := getPostId()
	if err != nil {
		return fmt.Errorf("getPostId: %v", err)
	}

	result := repo.db.First(&post, postID)
	if result.Error != nil {
		return fmt.Errorf("repo.db.First: %v", result.Error)
	}

	fmt.Print("Enter comment Content: ")
	fmt.Scanln(&content)

	comment := model.Comment{
		Content: content,
		PostID:  post.ID,
	}

	result = repo.db.Create(&comment)
	if result.Error != nil {
		return fmt.Errorf("repo.db.Create: %v", result.Error)
	}

	fmt.Printf("Comment created successfully!\n\n")
	return nil
}

func (repo Repository) Read() error {
	var opt int
	fmt.Print(
		"  1. Read posts\n",
		"  2. Read comments on post\n",
		"Choose action: ",
	)

	if _, err := fmt.Scanf("%d\n", &opt); err != nil {
		var temp string
		fmt.Scanln(&temp)
	}

	switch opt {
	case 1:
		if err := repo.readPosts(); err != nil {
			return fmt.Errorf("repo.readPosts: %v", err)
		}
	case 2:
		if err := repo.readComments(); err != nil {
			return fmt.Errorf("repo.readComments: %v", err)
		}
	default:
		fmt.Printf("Incorrect input\n\n")
	}
	return nil
}

func (repo Repository) readPosts() error {
	var posts []model.Post
	result := repo.db.Find(&posts)
	if result.Error != nil {
		return fmt.Errorf("repo.db.First: %v", result.Error)
	}

	printPosts(posts)
	return nil
}

func (repo Repository) readComments() error {
	var (
		post     model.Post
		comments []model.Comment
	)

	postID, err := getPostId()
	if err != nil {
		return fmt.Errorf("getPostId: %v", err)
	}
	result := repo.db.First(&post, postID)
	if result.Error != nil {
		return fmt.Errorf("repo.db.First: %v", result.Error)
	}
	result = repo.db.Where("post_id = ?", postID).Find(&comments)
	if result.Error != nil {
		return fmt.Errorf("repo.db.Where.Find: %v", result.Error)
	}

	fmt.Printf(
		"Title: %s\nContent: %s\nComments:\n",
		post.Title,
		post.Content,
	)

	printComments(comments)

	return nil
}

func printPosts(posts []model.Post) {
	for _, post := range posts {
		fmt.Printf(
			"PostID %d\nTitle: %s\nContent: %s\n\n",
			post.ID,
			post.Title,
			post.Content,
		)
	}
}

func printComments(comments []model.Comment) {
	for _, comment := range comments {
		fmt.Printf(
			"CommentID %d\nContent: %s\n\n",
			comment.ID,
			comment.Content,
		)
	}
}

func (repo Repository) Update() error {
	var opt int
	fmt.Print(
		"  1. Update post\n",
		"  2. Update comment\n",
		"Choose action: ",
	)

	if _, err := fmt.Scanf("%d\n", &opt); err != nil {
		var temp string
		fmt.Scanln(&temp)
	}

	switch opt {
	case 1:
		if err := repo.updatePost(); err != nil {
			return fmt.Errorf("repo.updatePost: %v", err)
		}
	case 2:
		if err := repo.updateComment(); err != nil {
			return fmt.Errorf("repo.updatePost: %v", err)
		}
	default:
		fmt.Printf("Incorrect input\n\n")
	}
	return nil
}

func (repo Repository) updatePost() error {
	var (
		post    model.Post
		title   string
		content string
	)

	postID, err := getPostId()
	if err != nil {
		return fmt.Errorf("getPostId: %v", err)
	}
	repo.db.First(&post, postID)

	fmt.Print("Enter new post Title: ")
	fmt.Scanln(&title)
	fmt.Print("Enter new post Content: ")
	fmt.Scanln(&content)

	result := repo.db.Model(&post).Updates(model.Post{Title: title, Content: content})
	if result.Error != nil {
		return fmt.Errorf("repo.db.Model.Updates: %v", result.Error)
	}

	fmt.Printf("Post updated successfully!\n\n")
	return nil
}

func (repo Repository) updateComment() error {
	var (
		comment model.Comment
		content string
	)

	commentID, err := getCommentId()
	if err != nil {
		return fmt.Errorf("getCommentId: %v", err)
	}
	result := repo.db.First(&comment, commentID)
	if result.Error != nil {
		return fmt.Errorf("repo.db.First: %v", result.Error)
	}

	fmt.Print("Enter new comment Content: ")
	fmt.Scanln(&content)

	result = repo.db.Model(&comment).Updates(model.Comment{Content: content})
	if result.Error != nil {
		return fmt.Errorf("repo.db.Model.Updates: %v", result.Error)
	}

	fmt.Printf("Comment updated successfully!\n\n")
	return nil
}

func (repo Repository) Delete() error {
	var opt int
	fmt.Print(
		"  1. Delete post\n",
		"  2. Delete comment\n",
		"Choose action: ",
	)

	if _, err := fmt.Scanf("%d\n", &opt); err != nil {
		var temp string
		fmt.Scanln(&temp)
	}

	switch opt {
	case 1:
		if err := repo.deletePost(); err != nil {
			return fmt.Errorf("repo.deletePost: %v", err)
		}
	case 2:
		if err := repo.deleteComment(); err != nil {
			return fmt.Errorf("repo.deleteComment: %v", err)
		}
	default:
		fmt.Printf("Incorrect input\n\n")
	}
	return nil
}

func (repo Repository) deletePost() error {
	postID, err := getPostId()
	if err != nil {
		return fmt.Errorf("getPostId: %v", err)
	}
	result := repo.db.Delete(&model.Post{}, postID)
	if result.Error != nil {
		return fmt.Errorf("repo.db.Delete: %v", result.Error)
	}

	fmt.Printf("Post deleted successfully!\n\n")
	return nil
}

func (repo Repository) deleteComment() error {
	commentID, err := getCommentId()
	if err != nil {
		return fmt.Errorf("getCommentId: %v", err)
	}
	result := repo.db.Delete(&model.Comment{}, commentID)
	if result.Error != nil {
		return fmt.Errorf("repo.db.Delete: %v", result.Error)
	}

	fmt.Printf("Comment deleted successfully!\n\n")
	return nil
}

func getPostId() (postID int, err error) {
	fmt.Print("Enter post ID: ")

	if _, err = fmt.Scanf("%d\n", &postID); err != nil {
		var temp string
		fmt.Scanln(&temp)
		return -1, fmt.Errorf("fmt.Scanf: %v", err)
	}
	return
}

func getCommentId() (commentID int, err error) {
	fmt.Print("Enter comment ID: ")

	if _, err := fmt.Scanf("%d\n", &commentID); err != nil {
		var temp string
		fmt.Scanln(&temp)
		return -1, fmt.Errorf("fmt.Scanf: %v", err)
	}
	return
}
