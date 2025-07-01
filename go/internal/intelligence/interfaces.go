package intelligence

import (
	"context"
	"github.com/aws-research-wizard/go/internal/data"
)

// DomainPackLoaderInterface defines the contract for loading domain packs
type DomainPackLoaderInterface interface {
	LoadDomainPack(domainName string) (*DomainPackInfo, error)
	LoadAllDomainPacks() (map[string]*DomainPackInfo, error)
	GetAvailableDomains() ([]string, error)
	ValidateDomainPack(domainName string) error
	ClearCache()
}

// RecommendationEngineInterface defines the contract for generating recommendations
type RecommendationEngineInterface interface {
	GenerateRecommendations(ctx context.Context, dataPath string) (*data.RecommendationResult, error)
}
