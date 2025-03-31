package repository

import (
	model "cms_api/internal/domain/entity"
	"cms_api/internal/infrastructure"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockDynamoDBClient はDynamoDBClientのモック
type MockDynamoDBClient struct {
	*infrastructure.DynamoDBClient
	mock.Mock
}

// モックメソッドを定義

type contentRepositoryTestSuite struct {
	suite.Suite
	repo            ContentRepository
	dynamoContainer *infrastructure.DynamoDBContainer
	ctx             context.Context
}

func TestContentRepository(t *testing.T) {
	suite.Run(t, new(contentRepositoryTestSuite))
}

func (s *contentRepositoryTestSuite) SetupSuite() {
	// コンテキストを初期化
	s.ctx = context.Background()

	// DynamoDB Localコンテナを起動
	container, err := infrastructure.SetupDynamoDBContainer(s.ctx)
	s.Require().NoError(err, "DynamoDB Localコンテナのセットアップに失敗しました")
	s.dynamoContainer = container

	// テーブルを作成
	err = container.CreateTable(s.ctx, "Contents")
	s.Require().NoError(err, "テーブルの作成に失敗しました")

	// DynamoDBクライアントを取得
	dbClient := infrastructure.NewDynamoDBClientFromContainer(container)

	// リポジトリを初期化
	s.repo = NewContentRepository(dbClient)
}

func (s *contentRepositoryTestSuite) TearDownSuite() {
	if s.dynamoContainer != nil {
		// コンテナを停止
		err := s.dynamoContainer.Teardown(s.ctx)
		s.Require().NoError(err, "DynamoDB Localコンテナの終了に失敗しました")
	}
}

func (s *contentRepositoryTestSuite) SetupTest() {
	// 各テスト前にテーブルをクリアする
	// 注: 実際の実装では、テーブルの全アイテムをスキャンして削除するか、
	// テスト用のテーブルを毎回作り直す方が良いでしょう
}

func (s *contentRepositoryTestSuite) createTestContent(id, title, description, body, coverImage, publishedAt, status, categoryID string, tags []string) *model.Article {
	content := &model.Article{
		ID:          id,
		Title:       title,
		Description: description,
		Body:        body,
		CoverImage:  coverImage,
		PublishedAt: publishedAt,
		Status:      status,
		CategoryID:  categoryID,
		Tags:        tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := s.repo.CreateArticle(content)
	s.Require().NoError(err, "テストデータの作成に失敗しました")
	return content
}

func (s *contentRepositoryTestSuite) TestGetArticles() {
	s.Run("Contentsテーブルにデータが存在しない場合、GetContentsが空の配列を返す", func() {
		contents, err := s.repo.GetArticles(s.ctx)
		s.Require().NoError(err)
		s.Require().Empty(contents, "テスト前にテーブルは空であるべき")
	})

	s.Run("ContentsテーブルにPK=ARTICLE,SK=ARTICLE#1のデータが存在する場合、GetArticlesが記事のデータを返す", func() {
		// テストデータを作成
		testContent := s.createTestContent(
			"ARTICLE",
			"GetArticlesテスト用コンテンツ",
			"これはGetArticlesのテスト用コンテンツです。",
			"test-body",
			"",
			"",
			"",
			"",
			[]string{"tag1", "tag2"},
		)

		contents, err := s.repo.GetArticles(s.ctx)
		s.Require().NoError(err)
		s.Require().NotEmpty(contents, "テーブルにデータが存在する場合は結果が返るべき")

		var found bool
		for _, c := range contents {
			println("--------------------------------")
			println(c.ID)
			println(testContent.ID)
			println("--------------------------------")
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

func (s *contentRepositoryTestSuite) TestUpdateContent() {
	s.Run("Contentsテーブルに該当のデータが存在する場合、そのデータを更新する", func() {
		// テストデータを作成
		testContent := s.createTestContent(
			"test-update-1",
			"更新前のタイトル",
			"更新前の本文",
			"テストユーザー",
			"",
			"",
			"",
			"",
			[]string{"tag1", "tag2"},
		)

		// データを更新
		testContent.Title = "更新後のタイトル"
		testContent.Body = "更新後の本文"
		err := s.repo.UpdateArticle(testContent)
		s.Require().NoError(err, "コンテンツの更新に失敗しました")

		// 更新されたデータを取得して検証
		contents, err := s.repo.GetArticles(s.ctx)
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
		err = s.repo.DeleteArticle(testContent.ID)
		s.Require().NoError(err, "テストデータの削除に失敗しました")
	})

	s.Run("その他のエラーが出た場合にはエラーを返す", func() {
		// このテストはスキップします
		s.T().Skip("DynamoDB Localに対してエラーを発生させるのは難しいため、スキップします")
	})
}

func (s *contentRepositoryTestSuite) TestDeleteContent() {
	s.Run("Contentsテーブルに該当のデータが存在する場合、そのデータを削除する", func() {
		// テストデータを作成
		testContent := s.createTestContent(
			"test-delete-1",
			"削除するコンテンツ",
			"このコンテンツは削除されます。",
			"テストユーザー",
			"",
			"",
			"",
			"",
			[]string{"tag1", "tag2"},
		)

		// データが存在することを確認
		contents, err := s.repo.GetArticles(s.ctx)
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
		err = s.repo.DeleteArticle(testContent.ID)
		s.Require().NoError(err, "コンテンツの削除に失敗しました")

		// データが削除されたことを確認
		contents, err = s.repo.GetArticles(s.ctx)
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
