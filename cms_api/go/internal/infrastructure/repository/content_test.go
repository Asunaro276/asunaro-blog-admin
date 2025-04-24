package repository

import (
	model "cms_api/internal/domain/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockDynamoDBClient はDynamoDBClientのモック
type MockDynamoDBClient struct {
	*dynamoDBClient
	mock.Mock
}

func TestContentRepository(t *testing.T) {
	suite.Run(t, new(dynamodbTestcontainersTestSuite))
}

func (s *dynamodbTestcontainersTestSuite) createTestContent(id, title, description, body, coverImage string, publishedAt time.Time, status string, categoryID string, tags []string) *model.Article {
	content := &model.Article{
		ID:          id,
		Title:       title,
		Description: description,
		Body:        body,
		CoverImage:  coverImage,
		Status:      status,
		CategoryID:  categoryID,
		Tags:        tags,
		PublishedAt: publishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := s.contentRepository.CreateArticle(content)
	s.Require().NoError(err, "テストデータの作成に失敗しました")
	return content
}

func (s *dynamodbTestcontainersTestSuite) TestGetArticles() {
	s.Run("Contentsテーブルにデータが存在しない場合、GetContentsが空の配列を返す", func() {
		contents, err := s.contentRepository.GetArticles(s.ctx)
		s.Require().NoError(err)
		s.Require().Empty(contents, "テスト前にテーブルは空であるべき")
	})

	s.Run("ContentsテーブルにPK=ARTICLE,SK=ARTICLE#1のデータが存在する場合、GetArticlesが記事のデータを返す", func() {
		// テストデータを作成
		testContent := s.createTestContent(
			"001",
			"AWS入門",
			"AWSの基本を解説します。",
			"AWSの基本として、EC2の使い方を解説します。",
			"https://example.com/aws-image.jpg",
			time.Now(),
			"published",
			"001",
			[]string{"tag1", "tag2"},
		)

		contents, err := s.contentRepository.GetArticles(s.ctx)
		s.Require().NoError(err)
		s.Require().NotEmpty(contents, "テーブルにデータが存在する場合は結果が返るべき")

		var found bool
		for _, c := range contents {
			if c.ID == testContent.ID {
				found = true
				s.Equal(testContent.Title, c.Title)
				s.Equal(testContent.Body, c.Body)
				break
			}
		}
		s.True(found, "作成したテストデータが取得できるべき")
	})

	s.Run("その他のエラーが出た場合にはエラーを返す", func() {
		s.T().Skip("DynamoDB Localに対してエラーを発生させるのは難しいため、スキップします")
	})
}

func (s *dynamodbTestcontainersTestSuite) TestUpdateContent() {
	s.T().Skip("UpdateArticleの実装はまだ完了していないためスキップします")

	s.Run("Contentsテーブルに該当のデータが存在する場合、そのデータを更新する", func() {
		// テストデータを作成
		testContent := s.createTestContent(
			"001",
			"AWS入門",
			"AWSの基本を解説します。",
			"AWSの基本として、EC2の使い方を解説します。",
			"https://example.com/aws-image.jpg",
			time.Now(),
			"published",
			"001",
			[]string{"tag1", "tag2"},
		)

		// データを更新
		testContent.Title = "更新後のタイトル"
		testContent.Body = "更新後の本文"
		err := s.contentRepository.UpdateArticle(testContent)
		s.Require().NoError(err, "コンテンツの更新に失敗しました")

		// 更新されたデータを取得して検証
		contents, err := s.contentRepository.GetArticles(s.ctx)
		s.Require().NoError(err)

		var updatedContent *model.Article
		for _, c := range contents {
			if c.ID == testContent.ID {
				tmp := c
				updatedContent = &tmp
				break
			}
		}

		s.Require().NotNil(updatedContent, "更新したコンテンツが取得できるべき")
		s.Equal("更新後のタイトル", updatedContent.Title)
		s.Equal("更新後の本文", updatedContent.Body)

		// テスト後のクリーンアップ
		err = s.contentRepository.DeleteArticle(testContent.ID)
		s.Require().NoError(err, "テストデータの削除に失敗しました")
	})

	s.Run("その他のエラーが出た場合にはエラーを返す", func() {
		// このテストはスキップします
		s.T().Skip("DynamoDB Localに対してエラーを発生させるのは難しいため、スキップします")
	})
}

func (s *dynamodbTestcontainersTestSuite) TestDeleteContent() {
	s.T().Skip("UpdateArticleの実装はまだ完了していないためスキップします")

	s.Run("Contentsテーブルに該当のデータが存在する場合、そのデータを削除する", func() {
		// テストデータを作成
		testContent := s.createTestContent(
			"001",
			"AWS入門",
			"AWSの基本を解説します。",
			"AWSの基本として、EC2の使い方を解説します。",
			"https://example.com/aws-image.jpg",
			time.Now(),
			"published",
			"001",
			[]string{"tag1", "tag2"},
		)

		// データが存在することを確認
		contents, err := s.contentRepository.GetArticles(s.ctx)
		s.Require().NoError(err)

		var exists bool
		for _, c := range contents {
			if c.ID == testContent.ID {
				exists = true
				break
			}
		}
		s.True(exists, "削除前にデータが存在するべき")

		// データを削除
		err = s.contentRepository.DeleteArticle(testContent.ID)
		s.Require().NoError(err, "コンテンツの削除に失敗しました")

		// データが削除されたことを確認
		contents, err = s.contentRepository.GetArticles(s.ctx)
		s.Require().NoError(err)

		for _, c := range contents {
			s.NotEqual(testContent.ID, c.ID, "削除したコンテンツが存在しないべき")
		}
	})

	s.Run("その他のエラーが出た場合にはエラーを返す", func() {
		// このテストはスキップします
		s.T().Skip("DynamoDB Localに対してエラーを発生させるのは難しいため、スキップします")
	})
}
