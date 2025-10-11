package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"elasticgaze/internal/models"
)

type CollectionsRepository struct {
	db *sql.DB
}

func NewCollectionsRepository(db *sql.DB) *CollectionsRepository {
	return &CollectionsRepository{db: db}
}

// Collections CRUD operations

func (r *CollectionsRepository) CreateCollection(req *models.CreateCollectionRequest) (*models.Collection, error) {
	query := `
		INSERT INTO tbl_collections (name, description) 
		VALUES (?, ?)
		RETURNING id, name, description, created_at, updated_at`

	var collection models.Collection
	err := r.db.QueryRow(query, req.Name, req.Description).Scan(
		&collection.ID,
		&collection.Name,
		&collection.Description,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return &collection, nil
}

func (r *CollectionsRepository) GetCollectionByID(id int) (*models.Collection, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM tbl_collections WHERE id = ?`

	var collection models.Collection
	err := r.db.QueryRow(query, id).Scan(
		&collection.ID,
		&collection.Name,
		&collection.Description,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return &collection, nil
}

func (r *CollectionsRepository) GetAllCollections() ([]*models.Collection, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM tbl_collections ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}
	defer rows.Close()

	var collections []*models.Collection
	for rows.Next() {
		var collection models.Collection
		err := rows.Scan(
			&collection.ID,
			&collection.Name,
			&collection.Description,
			&collection.CreatedAt,
			&collection.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collection: %w", err)
		}
		collections = append(collections, &collection)
	}

	return collections, nil
}

func (r *CollectionsRepository) UpdateCollection(id int, req *models.UpdateCollectionRequest) (*models.Collection, error) {
	// Build dynamic query based on provided fields
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}

	if len(setParts) == 0 {
		return r.GetCollectionByID(id) // Nothing to update
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE tbl_collections SET %s WHERE id = ?",
		fmt.Sprintf("%s", setParts[0]))
	for i := 1; i < len(setParts); i++ {
		query = query[:len(query)-11] + ", " + setParts[i] + " WHERE id = ?"
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection: %w", err)
	}

	return r.GetCollectionByID(id)
}

func (r *CollectionsRepository) DeleteCollection(id int) error {
	query := `DELETE FROM tbl_collections WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("collection not found")
	}

	return nil
}

// Folders CRUD operations

func (r *CollectionsRepository) CreateFolder(req *models.CreateFolderRequest) (*models.Folder, error) {
	query := `
		INSERT INTO tbl_folders (name, parent_folder_id, collection_id) 
		VALUES (?, ?, ?)
		RETURNING id, name, parent_folder_id, collection_id, created_at, updated_at`

	var folder models.Folder
	err := r.db.QueryRow(query, req.Name, req.ParentFolderID, req.CollectionID).Scan(
		&folder.ID,
		&folder.Name,
		&folder.ParentFolderID,
		&folder.CollectionID,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return &folder, nil
}

func (r *CollectionsRepository) GetFolderByID(id int) (*models.Folder, error) {
	query := `SELECT id, name, parent_folder_id, collection_id, created_at, updated_at FROM tbl_folders WHERE id = ?`

	var folder models.Folder
	err := r.db.QueryRow(query, id).Scan(
		&folder.ID,
		&folder.Name,
		&folder.ParentFolderID,
		&folder.CollectionID,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("folder not found")
		}
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	return &folder, nil
}

func (r *CollectionsRepository) GetFoldersByCollectionID(collectionID int) ([]*models.Folder, error) {
	query := `SELECT id, name, parent_folder_id, collection_id, created_at, updated_at FROM tbl_folders WHERE collection_id = ? ORDER BY name`

	rows, err := r.db.Query(query, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}
	defer rows.Close()

	var folders []*models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.ParentFolderID,
			&folder.CollectionID,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}
		folders = append(folders, &folder)
	}

	return folders, nil
}

func (r *CollectionsRepository) UpdateFolder(id int, req *models.UpdateFolderRequest) (*models.Folder, error) {
	// Build dynamic query based on provided fields
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.ParentFolderID != nil {
		setParts = append(setParts, "parent_folder_id = ?")
		args = append(args, *req.ParentFolderID)
	}

	if len(setParts) == 0 {
		return r.GetFolderByID(id) // Nothing to update
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE tbl_folders SET %s WHERE id = ?",
		fmt.Sprintf("%s", setParts[0]))
	for i := 1; i < len(setParts); i++ {
		query = query[:len(query)-11] + ", " + setParts[i] + " WHERE id = ?"
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	return r.GetFolderByID(id)
}

func (r *CollectionsRepository) DeleteFolder(id int) error {
	query := `DELETE FROM tbl_folders WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("folder not found")
	}

	return nil
}

// Requests CRUD operations

func (r *CollectionsRepository) CreateRequest(req *models.CreateRequestRequest) (*models.Request, error) {
	query := `
		INSERT INTO tbl_requests (name, method, url, body, description, folder_id, collection_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id, name, method, url, body, description, folder_id, collection_id, created_at, updated_at`

	var request models.Request
	err := r.db.QueryRow(query, req.Name, req.Method, req.URL, req.Body, req.Description, req.FolderID, req.CollectionID).Scan(
		&request.ID,
		&request.Name,
		&request.Method,
		&request.URL,
		&request.Body,
		&request.Description,
		&request.FolderID,
		&request.CollectionID,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return &request, nil
}

func (r *CollectionsRepository) GetRequestByID(id int) (*models.Request, error) {
	query := `SELECT id, name, method, url, body, description, folder_id, collection_id, created_at, updated_at FROM tbl_requests WHERE id = ?`

	var request models.Request
	err := r.db.QueryRow(query, id).Scan(
		&request.ID,
		&request.Name,
		&request.Method,
		&request.URL,
		&request.Body,
		&request.Description,
		&request.FolderID,
		&request.CollectionID,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("request not found")
		}
		return nil, fmt.Errorf("failed to get request: %w", err)
	}

	return &request, nil
}

func (r *CollectionsRepository) GetRequestsByCollectionID(collectionID int) ([]*models.Request, error) {
	query := `SELECT id, name, method, url, body, description, folder_id, collection_id, created_at, updated_at FROM tbl_requests WHERE collection_id = ? ORDER BY name`

	rows, err := r.db.Query(query, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}
	defer rows.Close()

	var requests []*models.Request
	for rows.Next() {
		var request models.Request
		err := rows.Scan(
			&request.ID,
			&request.Name,
			&request.Method,
			&request.URL,
			&request.Body,
			&request.Description,
			&request.FolderID,
			&request.CollectionID,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}
		requests = append(requests, &request)
	}

	return requests, nil
}

func (r *CollectionsRepository) GetRequestsByFolderID(folderID int) ([]*models.Request, error) {
	query := `SELECT id, name, method, url, body, description, folder_id, collection_id, created_at, updated_at FROM tbl_requests WHERE folder_id = ? ORDER BY name`

	rows, err := r.db.Query(query, folderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}
	defer rows.Close()

	var requests []*models.Request
	for rows.Next() {
		var request models.Request
		err := rows.Scan(
			&request.ID,
			&request.Name,
			&request.Method,
			&request.URL,
			&request.Body,
			&request.Description,
			&request.FolderID,
			&request.CollectionID,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}
		requests = append(requests, &request)
	}

	return requests, nil
}

func (r *CollectionsRepository) UpdateRequest(id int, req *models.UpdateRequestRequest) (*models.Request, error) {
	// Build dynamic query based on provided fields
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Method != nil {
		setParts = append(setParts, "method = ?")
		args = append(args, *req.Method)
	}
	if req.URL != nil {
		setParts = append(setParts, "url = ?")
		args = append(args, *req.URL)
	}
	if req.Body != nil {
		setParts = append(setParts, "body = ?")
		args = append(args, *req.Body)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.FolderID != nil {
		if *req.FolderID == -1 {
			// Special case: -1 means set to NULL
			setParts = append(setParts, "folder_id = NULL")
		} else {
			setParts = append(setParts, "folder_id = ?")
			args = append(args, *req.FolderID)
		}
	}
	if req.CollectionID != nil {
		setParts = append(setParts, "collection_id = ?")
		args = append(args, *req.CollectionID)
	}

	if len(setParts) == 0 {
		return r.GetRequestByID(id) // Nothing to update
	}

	// Build the complete query with proper JOIN syntax
	query := fmt.Sprintf("UPDATE tbl_requests SET %s WHERE id = ?", strings.Join(setParts, ", "))
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update request: %w", err)
	}

	return r.GetRequestByID(id)
}

func (r *CollectionsRepository) DeleteRequest(id int) error {
	query := `DELETE FROM tbl_requests WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("request not found")
	}

	return nil
}

// Tree structure methods

func (r *CollectionsRepository) GetCollectionTree(collectionID int) (*models.CollectionTreeNode, error) {
	// Get collection details
	collection, err := r.GetCollectionByID(collectionID)
	if err != nil {
		return nil, err
	}

	// Create root node
	root := &models.CollectionTreeNode{
		ID:       collection.ID,
		Name:     collection.Name,
		Type:     "collection",
		Children: []*models.CollectionTreeNode{},
	}

	// Get all folders and requests for this collection
	folders, err := r.GetFoldersByCollectionID(collectionID)
	if err != nil {
		return nil, err
	}

	requests, err := r.GetRequestsByCollectionID(collectionID)
	if err != nil {
		return nil, err
	}

	// Build tree structure
	r.buildTreeStructure(root, folders, requests)

	return root, nil
}

func (r *CollectionsRepository) buildTreeStructure(parent *models.CollectionTreeNode, folders []*models.Folder, requests []*models.Request) {
	// Add folders that belong to this parent
	var parentFolderID *int
	if parent.Type == "folder" {
		parentFolderID = &parent.ID
	} else if parent.Type == "collection" {
		parentFolderID = nil // Root level folders
	}

	for _, folder := range folders {
		if (parentFolderID == nil && folder.ParentFolderID == nil) ||
			(parentFolderID != nil && folder.ParentFolderID != nil && *folder.ParentFolderID == *parentFolderID) {

			folderNode := &models.CollectionTreeNode{
				ID:       folder.ID,
				Name:     folder.Name,
				Type:     "folder",
				Children: []*models.CollectionTreeNode{},
			}

			parent.Children = append(parent.Children, folderNode)

			// Recursively add children to this folder
			r.buildTreeStructure(folderNode, folders, requests)
		}
	}

	// Add requests that belong to this parent
	for _, request := range requests {
		var belongsToParent bool

		if parent.Type == "collection" && request.FolderID == nil {
			belongsToParent = true // Root level request
		} else if parent.Type == "folder" && request.FolderID != nil && *request.FolderID == parent.ID {
			belongsToParent = true // Request in this folder
		}

		if belongsToParent {
			requestNode := &models.CollectionTreeNode{
				ID:          request.ID,
				Name:        request.Name,
				Type:        "request",
				Method:      &request.Method,
				URL:         &request.URL,
				Body:        request.Body,
				Description: request.Description,
			}

			parent.Children = append(parent.Children, requestNode)
		}
	}
}
