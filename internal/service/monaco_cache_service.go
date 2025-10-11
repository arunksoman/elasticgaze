package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type MonacoCacheService struct {
	cacheDir string
}

type CacheInfo struct {
	Exists    bool   `json:"exists"`
	Size      int64  `json:"size"`
	ModTime   string `json:"modTime"`
	CacheKey  string `json:"cacheKey"`
	CachePath string `json:"cachePath"`
	IsExpired bool   `json:"isExpired"`
}

func NewMonacoCacheService(appDataDir string) *MonacoCacheService {
	cacheDir := filepath.Join(appDataDir, "monaco-cache")
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		fmt.Printf("Warning: Failed to create Monaco cache directory: %v\n", err)
	} else {
		fmt.Printf("Monaco cache directory created/verified: %s\n", cacheDir)
	}

	return &MonacoCacheService{
		cacheDir: cacheDir,
	}
}

func (s *MonacoCacheService) GetCacheInfo(version string) (*CacheInfo, error) {
	cacheKey := s.generateCacheKey(version)
	cachePath := filepath.Join(s.cacheDir, cacheKey+".cache")

	info := &CacheInfo{
		CacheKey:  cacheKey,
		CachePath: cachePath,
		Exists:    false,
		IsExpired: false,
	}

	stat, err := os.Stat(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return info, nil
		}
		return nil, err
	}

	info.Exists = true
	info.Size = stat.Size()
	info.ModTime = stat.ModTime().Format(time.RFC3339)

	// Check if cache is expired (older than 7 days)
	if time.Since(stat.ModTime()) > 7*24*time.Hour {
		info.IsExpired = true
	}

	return info, nil
}

func (s *MonacoCacheService) WriteCache(version string, data []byte) error {
	cacheKey := s.generateCacheKey(version)
	cachePath := filepath.Join(s.cacheDir, cacheKey+".cache")
	
	fmt.Printf("Writing Monaco cache to: %s (size: %d bytes)\n", cachePath, len(data))
	
	// Write data to temporary file first
	tempPath := cachePath + ".tmp"
	err := os.WriteFile(tempPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cache data: %w", err)
	}
	
	// Atomically move temp file to final location
	err = os.Rename(tempPath, cachePath)
	if err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to move cache file: %w", err)
	}
	
	fmt.Printf("Monaco cache written successfully: %s\n", cachePath)
	return nil
}

func (s *MonacoCacheService) ReadCache(version string) ([]byte, error) {
	cacheKey := s.generateCacheKey(version)
	cachePath := filepath.Join(s.cacheDir, cacheKey+".cache")

	fmt.Printf("Reading Monaco cache from: %s\n", cachePath)
	
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	fmt.Printf("Monaco cache read successfully: %d bytes\n", len(data))
	return data, nil
}

func (s *MonacoCacheService) InvalidateCache(version string) error {
	cacheKey := s.generateCacheKey(version)
	cachePath := filepath.Join(s.cacheDir, cacheKey+".cache")

	err := os.Remove(cachePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return nil
}

func (s *MonacoCacheService) ClearAllCache() error {
	entries, err := os.ReadDir(s.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) == ".cache" {
			err := os.Remove(filepath.Join(s.cacheDir, entry.Name()))
			if err != nil {
				return fmt.Errorf("failed to remove cache file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

func (s *MonacoCacheService) GetCacheSize() (int64, error) {
	var totalSize int64

	err := filepath.Walk(s.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".cache" {
			totalSize += info.Size()
		}

		return nil
	})

	return totalSize, err
}

func (s *MonacoCacheService) generateCacheKey(version string) string {
	hasher := md5.New()
	io.WriteString(hasher, "monaco-editor-"+version)
	return hex.EncodeToString(hasher.Sum(nil))
}
