package service

import (
	"fmt"

	"elasticgaze/internal/models"
	"elasticgaze/internal/repository"
)

type CollectionsService struct {
	collectionsRepo *repository.CollectionsRepository
}

func NewCollectionsService(collectionsRepo *repository.CollectionsRepository) *CollectionsService {
	return &CollectionsService{
		collectionsRepo: collectionsRepo,
	}
}

// Collections business logic

func (s *CollectionsService) CreateCollection(req *models.CreateCollectionRequest) (*models.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return s.collectionsRepo.CreateCollection(req)
}

func (s *CollectionsService) GetCollectionByID(id int) (*models.Collection, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid collection ID")
	}

	return s.collectionsRepo.GetCollectionByID(id)
}

func (s *CollectionsService) GetAllCollections() ([]*models.Collection, error) {
	return s.collectionsRepo.GetAllCollections()
}

func (s *CollectionsService) UpdateCollection(id int, req *models.UpdateCollectionRequest) (*models.Collection, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid collection ID")
	}

	// Check if collection exists
	_, err := s.collectionsRepo.GetCollectionByID(id)
	if err != nil {
		return nil, err
	}

	return s.collectionsRepo.UpdateCollection(id, req)
}

func (s *CollectionsService) DeleteCollection(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid collection ID")
	}

	// Check if collection exists
	_, err := s.collectionsRepo.GetCollectionByID(id)
	if err != nil {
		return err
	}

	return s.collectionsRepo.DeleteCollection(id)
}

// Folders business logic

func (s *CollectionsService) CreateFolder(req *models.CreateFolderRequest) (*models.Folder, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Validate that collection exists
	_, err := s.collectionsRepo.GetCollectionByID(req.CollectionID)
	if err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}

	// If parent folder ID is provided, validate it exists
	if req.ParentFolderID != nil {
		_, err := s.collectionsRepo.GetFolderByID(*req.ParentFolderID)
		if err != nil {
			return nil, fmt.Errorf("parent folder not found: %w", err)
		}
	}

	return s.collectionsRepo.CreateFolder(req)
}

func (s *CollectionsService) GetFolderByID(id int) (*models.Folder, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid folder ID")
	}

	return s.collectionsRepo.GetFolderByID(id)
}

func (s *CollectionsService) GetFoldersByCollectionID(collectionID int) ([]*models.Folder, error) {
	if collectionID <= 0 {
		return nil, fmt.Errorf("invalid collection ID")
	}

	return s.collectionsRepo.GetFoldersByCollectionID(collectionID)
}

func (s *CollectionsService) UpdateFolder(id int, req *models.UpdateFolderRequest) (*models.Folder, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid folder ID")
	}

	// Check if folder exists
	_, err := s.collectionsRepo.GetFolderByID(id)
	if err != nil {
		return nil, err
	}

	// If parent folder ID is being changed, validate it exists and doesn't create a circular reference
	if req.ParentFolderID != nil {
		if *req.ParentFolderID == id {
			return nil, fmt.Errorf("folder cannot be its own parent")
		}

		if *req.ParentFolderID != 0 { // 0 means root level
			_, err := s.collectionsRepo.GetFolderByID(*req.ParentFolderID)
			if err != nil {
				return nil, fmt.Errorf("parent folder not found: %w", err)
			}

			// Check for circular reference
			if s.wouldCreateCircularReference(id, *req.ParentFolderID) {
				return nil, fmt.Errorf("operation would create circular reference")
			}
		}
	}

	return s.collectionsRepo.UpdateFolder(id, req)
}

func (s *CollectionsService) DeleteFolder(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid folder ID")
	}

	// Check if folder exists
	_, err := s.collectionsRepo.GetFolderByID(id)
	if err != nil {
		return err
	}

	return s.collectionsRepo.DeleteFolder(id)
}

// Requests business logic

func (s *CollectionsService) CreateRequest(req *models.CreateRequestRequest) (*models.Request, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Validate that collection exists
	_, err := s.collectionsRepo.GetCollectionByID(req.CollectionID)
	if err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}

	// If folder ID is provided, validate it exists and belongs to the same collection
	if req.FolderID != nil {
		folder, err := s.collectionsRepo.GetFolderByID(*req.FolderID)
		if err != nil {
			return nil, fmt.Errorf("folder not found: %w", err)
		}
		if folder.CollectionID != req.CollectionID {
			return nil, fmt.Errorf("folder does not belong to the specified collection")
		}
	}

	return s.collectionsRepo.CreateRequest(req)
}

func (s *CollectionsService) GetRequestByID(id int) (*models.Request, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid request ID")
	}

	return s.collectionsRepo.GetRequestByID(id)
}

func (s *CollectionsService) GetRequestsByCollectionID(collectionID int) ([]*models.Request, error) {
	if collectionID <= 0 {
		return nil, fmt.Errorf("invalid collection ID")
	}

	return s.collectionsRepo.GetRequestsByCollectionID(collectionID)
}

func (s *CollectionsService) GetRequestsByFolderID(folderID int) ([]*models.Request, error) {
	if folderID <= 0 {
		return nil, fmt.Errorf("invalid folder ID")
	}

	return s.collectionsRepo.GetRequestsByFolderID(folderID)
}

func (s *CollectionsService) UpdateRequest(id int, req *models.UpdateRequestRequest) (*models.Request, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid request ID")
	}

	// Check if request exists
	request, err := s.collectionsRepo.GetRequestByID(id)
	if err != nil {
		return nil, err
	}

	// If folder ID is being changed, validate it exists and belongs to the same collection
	if req.FolderID != nil {
		if *req.FolderID == -1 {
			// Special case: -1 means move to collection root (no folder)
			// No validation needed, just pass through to repository
		} else if *req.FolderID != 0 { // 0 means root level
			folder, err := s.collectionsRepo.GetFolderByID(*req.FolderID)
			if err != nil {
				return nil, fmt.Errorf("folder not found: %w", err)
			}
			if folder.CollectionID != request.CollectionID {
				return nil, fmt.Errorf("folder does not belong to the same collection")
			}
		}
	}

	return s.collectionsRepo.UpdateRequest(id, req)
}

func (s *CollectionsService) DeleteRequest(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid request ID")
	}

	// Check if request exists
	_, err := s.collectionsRepo.GetRequestByID(id)
	if err != nil {
		return err
	}

	return s.collectionsRepo.DeleteRequest(id)
}

// Tree structure methods

func (s *CollectionsService) GetCollectionTree(collectionID int) (*models.CollectionTreeNode, error) {
	if collectionID <= 0 {
		return nil, fmt.Errorf("invalid collection ID")
	}

	return s.collectionsRepo.GetCollectionTree(collectionID)
}

func (s *CollectionsService) GetAllCollectionTrees() ([]*models.CollectionTreeNode, error) {
	collections, err := s.collectionsRepo.GetAllCollections()
	if err != nil {
		return nil, err
	}

	var trees []*models.CollectionTreeNode
	for _, collection := range collections {
		tree, err := s.collectionsRepo.GetCollectionTree(collection.ID)
		if err != nil {
			// Log error but continue with other collections
			continue
		}
		trees = append(trees, tree)
	}

	return trees, nil
}

// Helper methods

func (s *CollectionsService) wouldCreateCircularReference(folderID int, parentFolderID int) bool {
	// Check if making parentFolderID the parent of folderID would create a circular reference
	currentParentID := parentFolderID

	for currentParentID != 0 {
		if currentParentID == folderID {
			return true // Circular reference detected
		}

		// Get the parent of the current parent
		parent, err := s.collectionsRepo.GetFolderByID(currentParentID)
		if err != nil || parent.ParentFolderID == nil {
			break
		}

		currentParentID = *parent.ParentFolderID
	}

	return false
}

// Business logic for creating a default collection if none exists
func (s *CollectionsService) EnsureDefaultCollection() (*models.Collection, error) {
	collections, err := s.GetAllCollections()
	if err != nil {
		return nil, err
	}

	// If no collections exist, create a default one
	if len(collections) == 0 {
		defaultDescription := "Default collection for REST requests"
		req := &models.CreateCollectionRequest{
			Name:        "My Requests",
			Description: &defaultDescription,
		}
		return s.CreateCollection(req)
	}

	// Return the first collection as default
	return collections[0], nil
}
