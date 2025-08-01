package seeder

import (
	"context"
	"log"
	"trading-alchemist/internal/domain/chat"
	"trading-alchemist/internal/infrastructure/database"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
)

type ProviderSeed struct {
	Name        string
	DisplayName string
	Models      []ModelSeed
}

type ModelSeed struct {
	Name           string
	DisplayName    string
	SupportsVision bool
	IsActive       bool
}

var seeds = []ProviderSeed{
	{
		Name:        "openai",
		DisplayName: "OpenAI",
		Models: []ModelSeed{
			{Name: "gpt-4o", DisplayName: "GPT-4o", SupportsVision: true, IsActive: true},
			{Name: "gpt-4o-mini", DisplayName: "GPT-4o Mini", SupportsVision: true, IsActive: true},
			{Name: "gpt-4.1", DisplayName: "GPT-4.1", SupportsVision: true, IsActive: true}, // Assuming GPT-4.1 supports vision based on sources
			{Name: "gpt-4.1-mini", DisplayName: "GPT-4.1 Mini", SupportsVision: true, IsActive: true},
			{Name: "o3", DisplayName: "OpenAI o3", SupportsVision: true, IsActive: true},
			{Name: "o3-pro", DisplayName: "OpenAI o3-pro", SupportsVision: true, IsActive: true},
			{Name: "o4-mini", DisplayName: "OpenAI o4-mini", SupportsVision: true, IsActive: true},
		},
	},
	{
		Name:        "google",
		DisplayName: "Google",
		Models: []ModelSeed{
			{Name: "gemini-2.5-pro", DisplayName: "Gemini 2.5 Pro", SupportsVision: true, IsActive: true},
			{Name: "gemini-2.5-flash", DisplayName: "Gemini 2.5 Flash", SupportsVision: true, IsActive: true},
			{Name: "gemini-2.5-flash-lite", DisplayName: "Gemini 2.5 Flash-Lite", SupportsVision: true, IsActive: true},
		},
	},
}

func Seed(dbService *database.Service) {
	log.Println("Seeding providers and models...")

	err := dbService.ExecuteInTx(context.Background(), func(repoProvider database.RepositoryProvider) error {
		providerRepo := repoProvider.Provider()
		modelRepo := repoProvider.Model()

		for _, pSeed := range seeds {
			existingProvider, err := providerRepo.GetByName(context.Background(), pSeed.Name)
			if err != nil && err != errors.ErrProviderNotFound {
				return err
			}

			var providerID uuid.UUID
			if existingProvider != nil {
				providerID = existingProvider.ID
				log.Printf("Provider '%s' already exists.", pSeed.DisplayName)
			} else {
				newProvider := &chat.Provider{
					Name:        pSeed.Name,
					DisplayName: pSeed.DisplayName,
					IsActive:    true,
				}
				createdProvider, err := providerRepo.Create(context.Background(), newProvider)
				if err != nil {
					return err
				}
				providerID = createdProvider.ID
				log.Printf("Provider '%s' created.", pSeed.DisplayName)
			}

			for _, mSeed := range pSeed.Models {
				_, err := modelRepo.GetModelByName(context.Background(), providerID, mSeed.Name)
				if err != nil && err != errors.ErrModelNotFound {
					return err
				}
				if err == nil {
					log.Printf("Model '%s' for provider '%s' already exists.", mSeed.DisplayName, pSeed.DisplayName)
					continue
				}

				newModel := &chat.Model{
					ProviderID:     providerID,
					Name:           mSeed.Name,
					DisplayName:    mSeed.DisplayName,
					SupportsVision: mSeed.SupportsVision,
					IsActive:       mSeed.IsActive,
				}
				_, err = modelRepo.CreateModel(context.Background(), newModel)
				if err != nil {
					return err
				}
				log.Printf("Model '%s' for provider '%s' created.", mSeed.DisplayName, pSeed.DisplayName)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Seeding completed.")
}