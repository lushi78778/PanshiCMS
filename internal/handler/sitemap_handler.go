// File: internal/handler/sitemap_handler.go
package handler

import (
	"fmt"

	"PanshiCMS/internal/database"
	"PanshiCMS/internal/model"
	"github.com/gin-gonic/gin"
)

// SitemapHandler 输出网站的 sitemap.xml。该 sitemap 包含所有已发布的新闻、服务和案例。
// 访问路径为 /sitemap.xml。
func SitemapHandler(c *gin.Context) {
	// 查询数据
	var news []model.NewsArticle
	database.DB.Where("is_published = ?", true).Find(&news)
	var services []model.Service
	database.DB.Find(&services)
	var cases []model.CaseStudy
	database.DB.Where("is_published = ?", true).Find(&cases)
	// 站点基地址
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	// 构建 XML
	type urlEntry struct {
		Loc     string
		Lastmod string
	}
	var entries []urlEntry
	for _, n := range news {
		entries = append(entries, urlEntry{
			Loc:     fmt.Sprintf("%s/news/%s", baseURL, n.Slug),
			Lastmod: n.PublishDate.Format("2006-01-02"),
		})
	}
	for _, s := range services {
		entries = append(entries, urlEntry{
			Loc:     fmt.Sprintf("%s/services/%s", baseURL, s.Slug),
			Lastmod: s.UpdatedAt.Format("2006-01-02"),
		})
	}
	for _, cs := range cases {
		entries = append(entries, urlEntry{
			Loc:     fmt.Sprintf("%s/cases/%s", baseURL, cs.Slug),
			Lastmod: cs.UpdatedAt.Format("2006-01-02"),
		})
	}
	// 输出 XML
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.Writer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	c.Writer.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for _, u := range entries {
		c.Writer.WriteString("  <url>\n")
		c.Writer.WriteString(fmt.Sprintf("    <loc>%s</loc>\n", u.Loc))
		c.Writer.WriteString(fmt.Sprintf("    <lastmod>%s</lastmod>\n", u.Lastmod))
		c.Writer.WriteString("  </url>\n")
	}
	c.Writer.WriteString("</urlset>")
}
