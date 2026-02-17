package domain

import "time"

type Item struct {
	ID        string    `json:"id"         bson:"_id"`
	Name      string    `json:"name"        bson:"name"`
	Content   string    `json:"content"     bson:"content"`
	CreatedAt time.Time `json:"created_at"  bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at"  bson:"updated_at"`
}

type ListParams struct {
	Page     int    `form:"page"      json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	SortBy   string `form:"sort_by"   json:"sort_by"`
	SortDir  string `form:"sort_dir"  json:"sort_dir"`
}

func (p *ListParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortDir == "" {
		p.SortDir = "desc"
	}
}

type ListResult struct {
	Items      []*Item `json:"items"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}
