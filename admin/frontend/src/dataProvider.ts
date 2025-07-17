import { DataProvider } from "@refinedev/core";

const API_URL = "http://localhost:8080"; // CMS APIのURL

export const dataProvider: DataProvider = {
  getList: async ({ resource, pagination, filters, sorters, meta }) => {
    if (resource === "blog_posts") {
      const response = await fetch(`${API_URL}/`);
      const data = await response.json();
      
      return {
        data: data.map((article: any) => ({
          id: article.id,
          title: article.title,
          content: article.body,
          category: { id: article.category_id },
          status: article.status,
          createdAt: article.created_at,
        })),
        total: data.length,
      };
    }
    
    // categories用の仮データ
    if (resource === "categories") {
      return {
        data: [
          { id: "1", title: "Technology" },
          { id: "2", title: "Lifestyle" },
          { id: "3", title: "Business" },
        ],
        total: 3,
      };
    }
    
    return { data: [], total: 0 };
  },

  getOne: async ({ resource, id, meta }) => {
    if (resource === "blog_posts") {
      // 現在のAPIには個別記事取得がないため、全体から検索
      const response = await fetch(`${API_URL}/`);
      const data = await response.json();
      const article = data.find((item: any) => item.id === id);
      
      if (!article) {
        throw new Error("Article not found");
      }
      
      return {
        data: {
          id: article.id,
          title: article.title,
          content: article.body,
          category: { id: article.category_id },
          status: article.status,
          createdAt: article.created_at,
        },
      };
    }
    
    // categories用の仮データ
    if (resource === "categories") {
      const categories = [
        { id: "1", title: "Technology" },
        { id: "2", title: "Lifestyle" },
        { id: "3", title: "Business" },
      ];
      const category = categories.find((cat) => cat.id === id);
      return { data: category || null };
    }
    
    throw new Error("Resource not found");
  },

  create: async ({ resource, variables, meta }) => {
    // 現在のAPIには作成機能がないため、モックレスポンス
    if (resource === "blog_posts") {
      return {
        data: {
          id: Date.now().toString(),
          ...variables,
        },
      };
    }
    throw new Error("Create not implemented");
  },

  update: async ({ resource, id, variables, meta }) => {
    // 現在のAPIには更新機能がないため、モックレスポンス
    if (resource === "blog_posts") {
      return {
        data: {
          id,
          ...variables,
        },
      };
    }
    throw new Error("Update not implemented");
  },

  deleteOne: async ({ resource, id, meta }) => {
    // 現在のAPIには削除機能がないため、モックレスポンス
    if (resource === "blog_posts") {
      return {
        data: { id },
      };
    }
    throw new Error("Delete not implemented");
  },

  getApiUrl: () => API_URL,
}; 
