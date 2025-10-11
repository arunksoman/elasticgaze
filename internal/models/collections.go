package models

// Collection represents a collection of REST requests
type Collection struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
	CreatedAt   string  `json:"created_at" db:"created_at"`
	UpdatedAt   string  `json:"updated_at" db:"updated_at"`
}

// Folder represents a folder within a collection
type Folder struct {
	ID             int    `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	ParentFolderID *int   `json:"parent_folder_id,omitempty" db:"parent_folder_id"`
	CollectionID   int    `json:"collection_id" db:"collection_id"`
	CreatedAt      string `json:"created_at" db:"created_at"`
	UpdatedAt      string `json:"updated_at" db:"updated_at"`
}

// Request represents a REST request within a collection
type Request struct {
	ID           int     `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	Method       string  `json:"method" db:"method"`
	URL          string  `json:"url" db:"url"`
	Body         *string `json:"body,omitempty" db:"body"`
	Description  *string `json:"description,omitempty" db:"description"`
	FolderID     *int    `json:"folder_id,omitempty" db:"folder_id"`
	CollectionID int     `json:"collection_id" db:"collection_id"`
	CreatedAt    string  `json:"created_at" db:"created_at"`
	UpdatedAt    string  `json:"updated_at" db:"updated_at"`
}

// CollectionTreeNode represents a tree node in the collections hierarchy
type CollectionTreeNode struct {
	ID          int                   `json:"id"`
	Name        string                `json:"name"`
	Type        string                `json:"type"` // "collection", "folder", or "request"
	Method      *string               `json:"method,omitempty"`
	URL         *string               `json:"url,omitempty"`
	Body        *string               `json:"body,omitempty"`
	Description *string               `json:"description,omitempty"`
	Children    []*CollectionTreeNode `json:"children,omitempty"`
}

// CreateCollectionRequest represents the request payload for creating a new collection
type CreateCollectionRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}

// UpdateCollectionRequest represents the request payload for updating an existing collection
type UpdateCollectionRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CreateFolderRequest represents the request payload for creating a new folder
type CreateFolderRequest struct {
	Name           string `json:"name" validate:"required"`
	ParentFolderID *int   `json:"parent_folder_id,omitempty"`
	CollectionID   int    `json:"collection_id" validate:"required"`
}

// UpdateFolderRequest represents the request payload for updating an existing folder
type UpdateFolderRequest struct {
	Name           *string `json:"name,omitempty"`
	ParentFolderID *int    `json:"parent_folder_id,omitempty"`
}

// CreateRequestRequest represents the request payload for creating a new request
type CreateRequestRequest struct {
	Name         string  `json:"name" validate:"required"`
	Method       string  `json:"method" validate:"required"`
	URL          string  `json:"url" validate:"required"`
	Body         *string `json:"body,omitempty"`
	Description  *string `json:"description,omitempty"`
	FolderID     *int    `json:"folder_id,omitempty"`
	CollectionID int     `json:"collection_id" validate:"required"`
}

// UpdateRequestRequest represents the request payload for updating an existing request
type UpdateRequestRequest struct {
	Name         *string `json:"name,omitempty"`
	Method       *string `json:"method,omitempty"`
	URL          *string `json:"url,omitempty"`
	Body         *string `json:"body,omitempty"`
	Description  *string `json:"description,omitempty"`
	FolderID     *int    `json:"folder_id,omitempty"`
	CollectionID *int    `json:"collection_id,omitempty"`
}

// Validation methods
func (c *CreateCollectionRequest) Validate() error {
	if c.Name == "" {
		return ErrCollectionNameRequired
	}
	return nil
}

func (f *CreateFolderRequest) Validate() error {
	if f.Name == "" {
		return ErrFolderNameRequired
	}
	if f.CollectionID <= 0 {
		return ErrCollectionIDRequired
	}
	return nil
}

func (r *CreateRequestRequest) Validate() error {
	if r.Name == "" {
		return ErrRequestNameRequired
	}
	if r.Method == "" {
		return ErrMethodRequired
	}
	if r.URL == "" {
		return ErrURLRequired
	}
	if r.CollectionID <= 0 {
		return ErrCollectionIDRequired
	}
	return nil
}

// Collection validation errors
var (
	ErrCollectionNameRequired = &ValidationError{Field: "name", Message: "collection name is required"}
	ErrFolderNameRequired     = &ValidationError{Field: "name", Message: "folder name is required"}
	ErrRequestNameRequired    = &ValidationError{Field: "name", Message: "request name is required"}
	ErrURLRequired            = &ValidationError{Field: "url", Message: "URL is required"}
	ErrCollectionIDRequired   = &ValidationError{Field: "collection_id", Message: "collection ID is required"}
)
