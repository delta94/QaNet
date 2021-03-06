package actions

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/shanmugharajk/qanet/models"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

const question = "question"
const answer = "answer"

type commentForm struct {
	Type    string
	Comment string
	PostID  int64
	ID      int64
}

func AddComment(c buffalo.Context) error {
	commentForm := &commentForm{}

	if err := c.Bind(commentForm); err != nil {
		return errors.WithStack(err)
	}

	paramPostID := c.Param("postID")

	postID, err := strconv.ParseInt(paramPostID, 10, 64)

	if err != nil {
		return errors.WithStack(err)
	}

	commentForm.PostID = postID

	if commentForm.Type == question {
		return saveQuestionComment(c, commentForm)
	}

	if commentForm.Type == answer {
		return saveAnswerComment(c, commentForm)
	}

	// Should not reach here.
	return errors.New("internal error occurred")
}

func UpdateComment(c buffalo.Context) error {
	commentForm := &commentForm{}

	if err := c.Bind(commentForm); err != nil {
		return errors.WithStack(err)
	}

	paramPostID := c.Param("postID")
	paramCommentID := c.Param("commentID")

	postID, err := strconv.ParseInt(paramPostID, 10, 64)

	if err != nil {
		return errors.WithStack(err)
	}

	commentID, err := strconv.ParseInt(paramCommentID, 10, 64)

	if err != nil {
		return errors.WithStack(err)
	}

	commentForm.ID = commentID
	commentForm.PostID = postID

	if commentForm.Type == question {
		return updateQuestionComment(c, commentForm)
	}

	if commentForm.Type == answer {
		return updateAnswerComment(c, commentForm)
	}

	// Should not reach here.
	return errors.New("internal error occurred")
}

func DeleteComment(c buffalo.Context) error {
	paramPostType := c.Param("type")
	paramCommentID := c.Param("commentID")

	commentID, err := strconv.ParseInt(paramCommentID, 10, 64)

	if err != nil {
		return errors.WithStack(err)
	}

	tx, _ := c.Value("tx").(*gorm.DB)

	if paramPostType == question {
		return deleteComment(c, tx, commentID, models.QuestionComment{})
	}

	if paramPostType == answer {
		return deleteComment(c, tx, commentID, models.AnswerComment{})
	}

	// Should not reach here.
	return errors.New("internal error occurred")
}

func deleteComment(c buffalo.Context, tx *gorm.DB, id int64, model interface{}) error {
	if e := models.DeleteById(tx, id, model); e != nil {
		return errors.WithStack(e)
	}

	return c.Render(202, r.String("Deleted successfully"))
}

func saveAnswerComment(c buffalo.Context, commentForm *commentForm) error {
	comment := &models.AnswerComment{}
	comment.Comment = commentForm.Comment
	comment.AnswerID = commentForm.PostID
	comment.CreatedBy = c.Value("userId").(string)
	comment.UpdatedBy = c.Value("userId").(string)

	tx, _ := c.Value("tx").(*gorm.DB)

	verrors, err := models.Add(tx, comment)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrors.HasAny() {
		c.Set("errorMessage", "Invalid data has been entered. Please check and try again.")
		return c.Render(400, r.Template("text/html", "shared/_error"))
	}

	return renderAnswerComment(c, comment)
}

func updateAnswerComment(c buffalo.Context, commentForm *commentForm) error {
	var err error

	comment := &models.AnswerComment{}
	comment.Comment = commentForm.Comment
	comment.ID = commentForm.ID
	comment.AnswerID = commentForm.PostID
	comment.UpdatedBy = c.Value("userId").(string)

	tx, _ := c.Value("tx").(*gorm.DB)

	_, err = models.UpdateById(tx, comment.ID, comment)
	if err != nil {
		return errors.WithStack(err)
	}

	if err != nil {
		c.Set("errorMessage", "Invalid data has been entered. Please check and try again.")
		return c.Render(400, r.Template("text/html", "shared/_error"))
	}

	newComment := new(models.AnswerComment)
	err = models.GetById(tx, comment.ID, newComment)

	if err != nil {
		return errors.WithStack(err)
	}

	return renderAnswerComment(c, newComment)
}

func saveQuestionComment(c buffalo.Context, commentForm *commentForm) error {
	comment := &models.QuestionComment{}
	comment.Comment = commentForm.Comment
	comment.QuestionID = commentForm.PostID
	comment.CreatedBy = c.Value("userId").(string)
	comment.UpdatedBy = c.Value("userId").(string)

	tx, _ := c.Value("tx").(*gorm.DB)

	verrors, err := models.Add(tx, comment)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrors.HasAny() {
		c.Set("errorMessage", "Invalid data has been entered. Please check and try again.")
		return c.Render(400, r.Template("text/html", "shared/_error"))
	}

	return renderQuestionComment(c, comment)
}

func updateQuestionComment(c buffalo.Context, commentForm *commentForm) error {
	var err error

	comment := &models.QuestionComment{}
	comment.Comment = commentForm.Comment
	comment.ID = commentForm.ID
	comment.QuestionID = commentForm.PostID
	comment.UpdatedBy = c.Value("userId").(string)

	tx, _ := c.Value("tx").(*gorm.DB)

	_, err = models.UpdateById(tx, comment.ID, comment)

	if err != nil {
		return errors.WithStack(err)
	}

	if err != nil {
		c.Set("errorMessage", "Invalid data has been entered. Please check and try again.")
		return c.Render(400, r.Template("text/html", "shared/_error"))
	}

	newComment := new(models.QuestionComment)
	err = models.GetById(tx, comment.ID, newComment)

	if err != nil {
		return errors.WithStack(err)
	}

	return renderQuestionComment(c, newComment)
}

func renderQuestionComment(c buffalo.Context, comment *models.QuestionComment) error {
	c.Set("type", question)
	c.Set("id", comment.ID)
	c.Set("closeVotes", comment.CloseVotes)
	c.Set("comment", comment.Comment)
	c.Set("author", comment.CreatedBy)
	c.Set("at", comment.CreatedAt)

	return c.Render(200, r.Template("text/html", "shared/_comment"))
}

func renderAnswerComment(c buffalo.Context, comment *models.AnswerComment) error {
	c.Set("type", answer)
	c.Set("id", comment.ID)
	c.Set("closeVotes", comment.CloseVotes)
	c.Set("comment", comment.Comment)
	c.Set("author", comment.CreatedBy)
	c.Set("at", comment.CreatedAt)

	return c.Render(200, r.Template("text/html", "shared/_comment"))
}
