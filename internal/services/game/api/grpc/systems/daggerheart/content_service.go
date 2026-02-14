package daggerheart

import (
	"context"
	"errors"
	"strings"

	pb "github.com/louisbranch/fracturing.space/api/gen/go/systems/daggerheart/v1"
	"github.com/louisbranch/fracturing.space/internal/services/game/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DaggerheartContentService implements the Daggerheart content gRPC API.
type DaggerheartContentService struct {
	pb.UnimplementedDaggerheartContentServiceServer
	stores Stores
}

// NewDaggerheartContentService creates a configured gRPC handler for content catalog APIs.
func NewDaggerheartContentService(stores Stores) *DaggerheartContentService {
	return &DaggerheartContentService{stores: stores}
}

// GetContentCatalog returns the entire Daggerheart content catalog.
func (s *DaggerheartContentService) GetContentCatalog(ctx context.Context, _ *pb.GetDaggerheartContentCatalogRequest) (*pb.GetDaggerheartContentCatalogResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	classes, err := store.ListDaggerheartClasses(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list classes: %v", err)
	}
	subclasses, err := store.ListDaggerheartSubclasses(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list subclasses: %v", err)
	}
	heritages, err := store.ListDaggerheartHeritages(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list heritages: %v", err)
	}
	experiences, err := store.ListDaggerheartExperiences(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list experiences: %v", err)
	}
	adversaries, err := store.ListDaggerheartAdversaryEntries(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list adversaries: %v", err)
	}
	beastforms, err := store.ListDaggerheartBeastforms(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list beastforms: %v", err)
	}
	companionExperiences, err := store.ListDaggerheartCompanionExperiences(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list companion experiences: %v", err)
	}
	lootEntries, err := store.ListDaggerheartLootEntries(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list loot entries: %v", err)
	}
	damageTypes, err := store.ListDaggerheartDamageTypes(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list damage types: %v", err)
	}
	domains, err := store.ListDaggerheartDomains(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list domains: %v", err)
	}
	domainCards, err := store.ListDaggerheartDomainCards(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list domain cards: %v", err)
	}
	weapons, err := store.ListDaggerheartWeapons(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list weapons: %v", err)
	}
	armor, err := store.ListDaggerheartArmor(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list armor: %v", err)
	}
	items, err := store.ListDaggerheartItems(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list items: %v", err)
	}
	environments, err := store.ListDaggerheartEnvironments(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list environments: %v", err)
	}

	return &pb.GetDaggerheartContentCatalogResponse{
		Catalog: &pb.DaggerheartContentCatalog{
			Classes:              toProtoDaggerheartClasses(classes),
			Subclasses:           toProtoDaggerheartSubclasses(subclasses),
			Heritages:            toProtoDaggerheartHeritages(heritages),
			Experiences:          toProtoDaggerheartExperiences(experiences),
			Adversaries:          toProtoDaggerheartAdversaryEntries(adversaries),
			Beastforms:           toProtoDaggerheartBeastforms(beastforms),
			CompanionExperiences: toProtoDaggerheartCompanionExperiences(companionExperiences),
			LootEntries:          toProtoDaggerheartLootEntries(lootEntries),
			DamageTypes:          toProtoDaggerheartDamageTypes(damageTypes),
			Domains:              toProtoDaggerheartDomains(domains),
			DomainCards:          toProtoDaggerheartDomainCards(domainCards),
			Weapons:              toProtoDaggerheartWeapons(weapons),
			Armor:                toProtoDaggerheartArmorList(armor),
			Items:                toProtoDaggerheartItems(items),
			Environments:         toProtoDaggerheartEnvironments(environments),
		},
	}, nil
}

// GetClass returns a single Daggerheart class.
func (s *DaggerheartContentService) GetClass(ctx context.Context, in *pb.GetDaggerheartClassRequest) (*pb.GetDaggerheartClassResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "class request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "class id is required")
	}

	class, err := store.GetDaggerheartClass(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get class", err)
	}

	return &pb.GetDaggerheartClassResponse{Class: toProtoDaggerheartClass(class)}, nil
}

// ListClasses returns all Daggerheart classes.
func (s *DaggerheartContentService) ListClasses(ctx context.Context, _ *pb.ListDaggerheartClassesRequest) (*pb.ListDaggerheartClassesResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	classes, err := store.ListDaggerheartClasses(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list classes: %v", err)
	}
	return &pb.ListDaggerheartClassesResponse{Classes: toProtoDaggerheartClasses(classes)}, nil
}

// GetSubclass returns a single Daggerheart subclass.
func (s *DaggerheartContentService) GetSubclass(ctx context.Context, in *pb.GetDaggerheartSubclassRequest) (*pb.GetDaggerheartSubclassResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "subclass request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "subclass id is required")
	}

	subclass, err := store.GetDaggerheartSubclass(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get subclass", err)
	}

	return &pb.GetDaggerheartSubclassResponse{Subclass: toProtoDaggerheartSubclass(subclass)}, nil
}

// ListSubclasses returns all Daggerheart subclasses.
func (s *DaggerheartContentService) ListSubclasses(ctx context.Context, _ *pb.ListDaggerheartSubclassesRequest) (*pb.ListDaggerheartSubclassesResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	subclasses, err := store.ListDaggerheartSubclasses(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list subclasses: %v", err)
	}
	return &pb.ListDaggerheartSubclassesResponse{Subclasses: toProtoDaggerheartSubclasses(subclasses)}, nil
}

// GetHeritage returns a single Daggerheart heritage.
func (s *DaggerheartContentService) GetHeritage(ctx context.Context, in *pb.GetDaggerheartHeritageRequest) (*pb.GetDaggerheartHeritageResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "heritage request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "heritage id is required")
	}

	heritage, err := store.GetDaggerheartHeritage(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get heritage", err)
	}

	return &pb.GetDaggerheartHeritageResponse{Heritage: toProtoDaggerheartHeritage(heritage)}, nil
}

// ListHeritages returns all Daggerheart heritages.
func (s *DaggerheartContentService) ListHeritages(ctx context.Context, _ *pb.ListDaggerheartHeritagesRequest) (*pb.ListDaggerheartHeritagesResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	heritages, err := store.ListDaggerheartHeritages(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list heritages: %v", err)
	}
	return &pb.ListDaggerheartHeritagesResponse{Heritages: toProtoDaggerheartHeritages(heritages)}, nil
}

// GetExperience returns a single Daggerheart experience.
func (s *DaggerheartContentService) GetExperience(ctx context.Context, in *pb.GetDaggerheartExperienceRequest) (*pb.GetDaggerheartExperienceResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "experience request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "experience id is required")
	}

	experience, err := store.GetDaggerheartExperience(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get experience", err)
	}

	return &pb.GetDaggerheartExperienceResponse{Experience: toProtoDaggerheartExperience(experience)}, nil
}

// ListExperiences returns all Daggerheart experiences.
func (s *DaggerheartContentService) ListExperiences(ctx context.Context, _ *pb.ListDaggerheartExperiencesRequest) (*pb.ListDaggerheartExperiencesResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	experiences, err := store.ListDaggerheartExperiences(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list experiences: %v", err)
	}
	return &pb.ListDaggerheartExperiencesResponse{Experiences: toProtoDaggerheartExperiences(experiences)}, nil
}

// GetAdversary returns a single Daggerheart adversary catalog entry.
func (s *DaggerheartContentService) GetAdversary(ctx context.Context, in *pb.GetDaggerheartAdversaryRequest) (*pb.GetDaggerheartAdversaryResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "adversary request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "adversary id is required")
	}

	adversary, err := store.GetDaggerheartAdversaryEntry(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get adversary", err)
	}

	return &pb.GetDaggerheartAdversaryResponse{Adversary: toProtoDaggerheartAdversaryEntry(adversary)}, nil
}

// ListAdversaries returns all Daggerheart adversary catalog entries.
func (s *DaggerheartContentService) ListAdversaries(ctx context.Context, _ *pb.ListDaggerheartAdversariesRequest) (*pb.ListDaggerheartAdversariesResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	adversaries, err := store.ListDaggerheartAdversaryEntries(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list adversaries: %v", err)
	}
	return &pb.ListDaggerheartAdversariesResponse{Adversaries: toProtoDaggerheartAdversaryEntries(adversaries)}, nil
}

// GetBeastform returns a single Daggerheart beastform catalog entry.
func (s *DaggerheartContentService) GetBeastform(ctx context.Context, in *pb.GetDaggerheartBeastformRequest) (*pb.GetDaggerheartBeastformResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "beastform request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "beastform id is required")
	}

	beastform, err := store.GetDaggerheartBeastform(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get beastform", err)
	}

	return &pb.GetDaggerheartBeastformResponse{Beastform: toProtoDaggerheartBeastform(beastform)}, nil
}

// ListBeastforms returns all Daggerheart beastform catalog entries.
func (s *DaggerheartContentService) ListBeastforms(ctx context.Context, _ *pb.ListDaggerheartBeastformsRequest) (*pb.ListDaggerheartBeastformsResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	beastforms, err := store.ListDaggerheartBeastforms(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list beastforms: %v", err)
	}
	return &pb.ListDaggerheartBeastformsResponse{Beastforms: toProtoDaggerheartBeastforms(beastforms)}, nil
}

// GetCompanionExperience returns a single Daggerheart companion experience catalog entry.
func (s *DaggerheartContentService) GetCompanionExperience(ctx context.Context, in *pb.GetDaggerheartCompanionExperienceRequest) (*pb.GetDaggerheartCompanionExperienceResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "companion experience request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "companion experience id is required")
	}

	experience, err := store.GetDaggerheartCompanionExperience(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get companion experience", err)
	}

	return &pb.GetDaggerheartCompanionExperienceResponse{Experience: toProtoDaggerheartCompanionExperience(experience)}, nil
}

// ListCompanionExperiences returns all Daggerheart companion experience catalog entries.
func (s *DaggerheartContentService) ListCompanionExperiences(ctx context.Context, _ *pb.ListDaggerheartCompanionExperiencesRequest) (*pb.ListDaggerheartCompanionExperiencesResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	experiences, err := store.ListDaggerheartCompanionExperiences(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list companion experiences: %v", err)
	}
	return &pb.ListDaggerheartCompanionExperiencesResponse{Experiences: toProtoDaggerheartCompanionExperiences(experiences)}, nil
}

// GetLootEntry returns a single Daggerheart loot catalog entry.
func (s *DaggerheartContentService) GetLootEntry(ctx context.Context, in *pb.GetDaggerheartLootEntryRequest) (*pb.GetDaggerheartLootEntryResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "loot entry request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "loot entry id is required")
	}

	entry, err := store.GetDaggerheartLootEntry(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get loot entry", err)
	}

	return &pb.GetDaggerheartLootEntryResponse{Entry: toProtoDaggerheartLootEntry(entry)}, nil
}

// ListLootEntries returns all Daggerheart loot catalog entries.
func (s *DaggerheartContentService) ListLootEntries(ctx context.Context, _ *pb.ListDaggerheartLootEntriesRequest) (*pb.ListDaggerheartLootEntriesResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	entries, err := store.ListDaggerheartLootEntries(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list loot entries: %v", err)
	}
	return &pb.ListDaggerheartLootEntriesResponse{Entries: toProtoDaggerheartLootEntries(entries)}, nil
}

// GetDamageType returns a single Daggerheart damage type catalog entry.
func (s *DaggerheartContentService) GetDamageType(ctx context.Context, in *pb.GetDaggerheartDamageTypeRequest) (*pb.GetDaggerheartDamageTypeResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "damage type request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "damage type id is required")
	}

	entry, err := store.GetDaggerheartDamageType(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get damage type", err)
	}

	return &pb.GetDaggerheartDamageTypeResponse{DamageType: toProtoDaggerheartDamageType(entry)}, nil
}

// ListDamageTypes returns all Daggerheart damage type catalog entries.
func (s *DaggerheartContentService) ListDamageTypes(ctx context.Context, _ *pb.ListDaggerheartDamageTypesRequest) (*pb.ListDaggerheartDamageTypesResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	entries, err := store.ListDaggerheartDamageTypes(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list damage types: %v", err)
	}
	return &pb.ListDaggerheartDamageTypesResponse{DamageTypes: toProtoDaggerheartDamageTypes(entries)}, nil
}

// GetDomain returns a single Daggerheart domain.
func (s *DaggerheartContentService) GetDomain(ctx context.Context, in *pb.GetDaggerheartDomainRequest) (*pb.GetDaggerheartDomainResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "domain request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "domain id is required")
	}

	domain, err := store.GetDaggerheartDomain(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get domain", err)
	}

	return &pb.GetDaggerheartDomainResponse{Domain: toProtoDaggerheartDomain(domain)}, nil
}

// ListDomains returns all Daggerheart domains.
func (s *DaggerheartContentService) ListDomains(ctx context.Context, _ *pb.ListDaggerheartDomainsRequest) (*pb.ListDaggerheartDomainsResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	domains, err := store.ListDaggerheartDomains(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list domains: %v", err)
	}
	return &pb.ListDaggerheartDomainsResponse{Domains: toProtoDaggerheartDomains(domains)}, nil
}

// GetDomainCard returns a single Daggerheart domain card.
func (s *DaggerheartContentService) GetDomainCard(ctx context.Context, in *pb.GetDaggerheartDomainCardRequest) (*pb.GetDaggerheartDomainCardResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "domain card request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "domain card id is required")
	}

	card, err := store.GetDaggerheartDomainCard(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get domain card", err)
	}

	return &pb.GetDaggerheartDomainCardResponse{DomainCard: toProtoDaggerheartDomainCard(card)}, nil
}

// ListDomainCards returns Daggerheart domain cards, optionally filtered by domain.
func (s *DaggerheartContentService) ListDomainCards(ctx context.Context, in *pb.ListDaggerheartDomainCardsRequest) (*pb.ListDaggerheartDomainCardsResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	var cards []storage.DaggerheartDomainCard
	if in != nil && strings.TrimSpace(in.GetDomainId()) != "" {
		cards, err = store.ListDaggerheartDomainCardsByDomain(ctx, in.GetDomainId())
	} else {
		cards, err = store.ListDaggerheartDomainCards(ctx)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list domain cards: %v", err)
	}

	return &pb.ListDaggerheartDomainCardsResponse{DomainCards: toProtoDaggerheartDomainCards(cards)}, nil
}

// GetWeapon returns a single Daggerheart weapon.
func (s *DaggerheartContentService) GetWeapon(ctx context.Context, in *pb.GetDaggerheartWeaponRequest) (*pb.GetDaggerheartWeaponResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "weapon request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "weapon id is required")
	}

	weapon, err := store.GetDaggerheartWeapon(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get weapon", err)
	}

	return &pb.GetDaggerheartWeaponResponse{Weapon: toProtoDaggerheartWeapon(weapon)}, nil
}

// ListWeapons returns all Daggerheart weapons.
func (s *DaggerheartContentService) ListWeapons(ctx context.Context, _ *pb.ListDaggerheartWeaponsRequest) (*pb.ListDaggerheartWeaponsResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	weapons, err := store.ListDaggerheartWeapons(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list weapons: %v", err)
	}
	return &pb.ListDaggerheartWeaponsResponse{Weapons: toProtoDaggerheartWeapons(weapons)}, nil
}

// GetArmor returns a single Daggerheart armor entry.
func (s *DaggerheartContentService) GetArmor(ctx context.Context, in *pb.GetDaggerheartArmorRequest) (*pb.GetDaggerheartArmorResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "armor request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "armor id is required")
	}

	armor, err := store.GetDaggerheartArmor(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get armor", err)
	}

	return &pb.GetDaggerheartArmorResponse{Armor: toProtoDaggerheartArmor(armor)}, nil
}

// ListArmor returns all Daggerheart armor entries.
func (s *DaggerheartContentService) ListArmor(ctx context.Context, _ *pb.ListDaggerheartArmorRequest) (*pb.ListDaggerheartArmorResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	armor, err := store.ListDaggerheartArmor(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list armor: %v", err)
	}
	return &pb.ListDaggerheartArmorResponse{Armor: toProtoDaggerheartArmorList(armor)}, nil
}

// GetItem returns a single Daggerheart item.
func (s *DaggerheartContentService) GetItem(ctx context.Context, in *pb.GetDaggerheartItemRequest) (*pb.GetDaggerheartItemResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "item request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "item id is required")
	}

	item, err := store.GetDaggerheartItem(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get item", err)
	}

	return &pb.GetDaggerheartItemResponse{Item: toProtoDaggerheartItem(item)}, nil
}

// ListItems returns all Daggerheart items.
func (s *DaggerheartContentService) ListItems(ctx context.Context, _ *pb.ListDaggerheartItemsRequest) (*pb.ListDaggerheartItemsResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	items, err := store.ListDaggerheartItems(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list items: %v", err)
	}
	return &pb.ListDaggerheartItemsResponse{Items: toProtoDaggerheartItems(items)}, nil
}

// GetEnvironment returns a single Daggerheart environment.
func (s *DaggerheartContentService) GetEnvironment(ctx context.Context, in *pb.GetDaggerheartEnvironmentRequest) (*pb.GetDaggerheartEnvironmentResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "environment request is required")
	}
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "environment id is required")
	}

	env, err := store.GetDaggerheartEnvironment(ctx, in.GetId())
	if err != nil {
		return nil, mapContentErr("get environment", err)
	}

	return &pb.GetDaggerheartEnvironmentResponse{Environment: toProtoDaggerheartEnvironment(env)}, nil
}

// ListEnvironments returns all Daggerheart environments.
func (s *DaggerheartContentService) ListEnvironments(ctx context.Context, _ *pb.ListDaggerheartEnvironmentsRequest) (*pb.ListDaggerheartEnvironmentsResponse, error) {
	store, err := s.contentStore()
	if err != nil {
		return nil, err
	}

	items, err := store.ListDaggerheartEnvironments(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list environments: %v", err)
	}
	return &pb.ListDaggerheartEnvironmentsResponse{Environments: toProtoDaggerheartEnvironments(items)}, nil
}

func (s *DaggerheartContentService) contentStore() (storage.DaggerheartContentStore, error) {
	if s == nil || s.stores.DaggerheartContent == nil {
		return nil, status.Error(codes.Internal, "content store is not configured")
	}
	return s.stores.DaggerheartContent, nil
}

func mapContentErr(action string, err error) error {
	if errors.Is(err, storage.ErrNotFound) {
		return status.Error(codes.NotFound, "content not found")
	}
	return status.Errorf(codes.Internal, "%s: %v", action, err)
}

func toProtoDaggerheartClass(class storage.DaggerheartClass) *pb.DaggerheartClass {
	return &pb.DaggerheartClass{
		Id:              class.ID,
		Name:            class.Name,
		StartingEvasion: int32(class.StartingEvasion),
		StartingHp:      int32(class.StartingHP),
		StartingItems:   append([]string{}, class.StartingItems...),
		Features:        toProtoDaggerheartFeatures(class.Features),
		HopeFeature:     toProtoDaggerheartHopeFeature(class.HopeFeature),
		DomainIds:       append([]string{}, class.DomainIDs...),
	}
}

func toProtoDaggerheartClasses(classes []storage.DaggerheartClass) []*pb.DaggerheartClass {
	items := make([]*pb.DaggerheartClass, 0, len(classes))
	for _, class := range classes {
		items = append(items, toProtoDaggerheartClass(class))
	}
	return items
}

func toProtoDaggerheartSubclass(subclass storage.DaggerheartSubclass) *pb.DaggerheartSubclass {
	return &pb.DaggerheartSubclass{
		Id:                     subclass.ID,
		Name:                   subclass.Name,
		SpellcastTrait:         subclass.SpellcastTrait,
		FoundationFeatures:     toProtoDaggerheartFeatures(subclass.FoundationFeatures),
		SpecializationFeatures: toProtoDaggerheartFeatures(subclass.SpecializationFeatures),
		MasteryFeatures:        toProtoDaggerheartFeatures(subclass.MasteryFeatures),
	}
}

func toProtoDaggerheartSubclasses(subclasses []storage.DaggerheartSubclass) []*pb.DaggerheartSubclass {
	items := make([]*pb.DaggerheartSubclass, 0, len(subclasses))
	for _, subclass := range subclasses {
		items = append(items, toProtoDaggerheartSubclass(subclass))
	}
	return items
}

func toProtoDaggerheartHeritage(heritage storage.DaggerheartHeritage) *pb.DaggerheartHeritage {
	return &pb.DaggerheartHeritage{
		Id:       heritage.ID,
		Name:     heritage.Name,
		Kind:     heritageKindToProto(heritage.Kind),
		Features: toProtoDaggerheartFeatures(heritage.Features),
	}
}

func toProtoDaggerheartExperience(experience storage.DaggerheartExperienceEntry) *pb.DaggerheartExperienceEntry {
	return &pb.DaggerheartExperienceEntry{
		Id:          experience.ID,
		Name:        experience.Name,
		Description: experience.Description,
	}
}

func toProtoDaggerheartExperiences(experiences []storage.DaggerheartExperienceEntry) []*pb.DaggerheartExperienceEntry {
	items := make([]*pb.DaggerheartExperienceEntry, 0, len(experiences))
	for _, experience := range experiences {
		items = append(items, toProtoDaggerheartExperience(experience))
	}
	return items
}

func toProtoDaggerheartAdversaryAttack(attack storage.DaggerheartAdversaryAttack) *pb.DaggerheartAdversaryAttack {
	return &pb.DaggerheartAdversaryAttack{
		Name:        attack.Name,
		Range:       attack.Range,
		DamageDice:  toProtoDaggerheartDamageDice(attack.DamageDice),
		DamageBonus: int32(attack.DamageBonus),
		DamageType:  damageTypeToProto(attack.DamageType),
	}
}

func toProtoDaggerheartAdversaryExperiences(experiences []storage.DaggerheartAdversaryExperience) []*pb.DaggerheartAdversaryExperience {
	items := make([]*pb.DaggerheartAdversaryExperience, 0, len(experiences))
	for _, experience := range experiences {
		items = append(items, &pb.DaggerheartAdversaryExperience{
			Name:     experience.Name,
			Modifier: int32(experience.Modifier),
		})
	}
	return items
}

func toProtoDaggerheartAdversaryFeatures(features []storage.DaggerheartAdversaryFeature) []*pb.DaggerheartAdversaryFeature {
	items := make([]*pb.DaggerheartAdversaryFeature, 0, len(features))
	for _, feature := range features {
		items = append(items, &pb.DaggerheartAdversaryFeature{
			Id:          feature.ID,
			Name:        feature.Name,
			Kind:        feature.Kind,
			Description: feature.Description,
			CostType:    feature.CostType,
			Cost:        int32(feature.Cost),
		})
	}
	return items
}

func toProtoDaggerheartAdversaryEntry(entry storage.DaggerheartAdversaryEntry) *pb.DaggerheartAdversaryEntry {
	return &pb.DaggerheartAdversaryEntry{
		Id:              entry.ID,
		Name:            entry.Name,
		Tier:            int32(entry.Tier),
		Role:            entry.Role,
		Description:     entry.Description,
		Motives:         entry.Motives,
		Difficulty:      int32(entry.Difficulty),
		MajorThreshold:  int32(entry.MajorThreshold),
		SevereThreshold: int32(entry.SevereThreshold),
		Hp:              int32(entry.HP),
		Stress:          int32(entry.Stress),
		Armor:           int32(entry.Armor),
		AttackModifier:  int32(entry.AttackModifier),
		StandardAttack:  toProtoDaggerheartAdversaryAttack(entry.StandardAttack),
		Experiences:     toProtoDaggerheartAdversaryExperiences(entry.Experiences),
		Features:        toProtoDaggerheartAdversaryFeatures(entry.Features),
	}
}

func toProtoDaggerheartAdversaryEntries(entries []storage.DaggerheartAdversaryEntry) []*pb.DaggerheartAdversaryEntry {
	items := make([]*pb.DaggerheartAdversaryEntry, 0, len(entries))
	for _, entry := range entries {
		items = append(items, toProtoDaggerheartAdversaryEntry(entry))
	}
	return items
}

func toProtoDaggerheartBeastformAttack(attack storage.DaggerheartBeastformAttack) *pb.DaggerheartBeastformAttack {
	return &pb.DaggerheartBeastformAttack{
		Range:       attack.Range,
		Trait:       attack.Trait,
		DamageDice:  toProtoDaggerheartDamageDice(attack.DamageDice),
		DamageBonus: int32(attack.DamageBonus),
		DamageType:  damageTypeToProto(attack.DamageType),
	}
}

func toProtoDaggerheartBeastformFeatures(features []storage.DaggerheartBeastformFeature) []*pb.DaggerheartBeastformFeature {
	items := make([]*pb.DaggerheartBeastformFeature, 0, len(features))
	for _, feature := range features {
		items = append(items, &pb.DaggerheartBeastformFeature{
			Id:          feature.ID,
			Name:        feature.Name,
			Description: feature.Description,
		})
	}
	return items
}

func toProtoDaggerheartBeastform(beastform storage.DaggerheartBeastformEntry) *pb.DaggerheartBeastformEntry {
	return &pb.DaggerheartBeastformEntry{
		Id:           beastform.ID,
		Name:         beastform.Name,
		Tier:         int32(beastform.Tier),
		Examples:     beastform.Examples,
		Trait:        beastform.Trait,
		TraitBonus:   int32(beastform.TraitBonus),
		EvasionBonus: int32(beastform.EvasionBonus),
		Attack:       toProtoDaggerheartBeastformAttack(beastform.Attack),
		Advantages:   append([]string{}, beastform.Advantages...),
		Features:     toProtoDaggerheartBeastformFeatures(beastform.Features),
	}
}

func toProtoDaggerheartBeastforms(beastforms []storage.DaggerheartBeastformEntry) []*pb.DaggerheartBeastformEntry {
	items := make([]*pb.DaggerheartBeastformEntry, 0, len(beastforms))
	for _, beastform := range beastforms {
		items = append(items, toProtoDaggerheartBeastform(beastform))
	}
	return items
}

func toProtoDaggerheartCompanionExperience(experience storage.DaggerheartCompanionExperienceEntry) *pb.DaggerheartCompanionExperienceEntry {
	return &pb.DaggerheartCompanionExperienceEntry{
		Id:          experience.ID,
		Name:        experience.Name,
		Description: experience.Description,
	}
}

func toProtoDaggerheartCompanionExperiences(experiences []storage.DaggerheartCompanionExperienceEntry) []*pb.DaggerheartCompanionExperienceEntry {
	items := make([]*pb.DaggerheartCompanionExperienceEntry, 0, len(experiences))
	for _, experience := range experiences {
		items = append(items, toProtoDaggerheartCompanionExperience(experience))
	}
	return items
}

func toProtoDaggerheartLootEntry(entry storage.DaggerheartLootEntry) *pb.DaggerheartLootEntry {
	return &pb.DaggerheartLootEntry{
		Id:          entry.ID,
		Name:        entry.Name,
		Roll:        int32(entry.Roll),
		Description: entry.Description,
	}
}

func toProtoDaggerheartLootEntries(entries []storage.DaggerheartLootEntry) []*pb.DaggerheartLootEntry {
	items := make([]*pb.DaggerheartLootEntry, 0, len(entries))
	for _, entry := range entries {
		items = append(items, toProtoDaggerheartLootEntry(entry))
	}
	return items
}

func toProtoDaggerheartDamageType(entry storage.DaggerheartDamageTypeEntry) *pb.DaggerheartDamageTypeEntry {
	return &pb.DaggerheartDamageTypeEntry{
		Id:          entry.ID,
		Name:        entry.Name,
		Description: entry.Description,
	}
}

func toProtoDaggerheartDamageTypes(entries []storage.DaggerheartDamageTypeEntry) []*pb.DaggerheartDamageTypeEntry {
	items := make([]*pb.DaggerheartDamageTypeEntry, 0, len(entries))
	for _, entry := range entries {
		items = append(items, toProtoDaggerheartDamageType(entry))
	}
	return items
}

func toProtoDaggerheartHeritages(heritages []storage.DaggerheartHeritage) []*pb.DaggerheartHeritage {
	items := make([]*pb.DaggerheartHeritage, 0, len(heritages))
	for _, heritage := range heritages {
		items = append(items, toProtoDaggerheartHeritage(heritage))
	}
	return items
}

func toProtoDaggerheartDomain(domain storage.DaggerheartDomain) *pb.DaggerheartDomain {
	return &pb.DaggerheartDomain{
		Id:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
	}
}

func toProtoDaggerheartDomains(domains []storage.DaggerheartDomain) []*pb.DaggerheartDomain {
	items := make([]*pb.DaggerheartDomain, 0, len(domains))
	for _, domain := range domains {
		items = append(items, toProtoDaggerheartDomain(domain))
	}
	return items
}

func toProtoDaggerheartDomainCard(card storage.DaggerheartDomainCard) *pb.DaggerheartDomainCard {
	return &pb.DaggerheartDomainCard{
		Id:          card.ID,
		Name:        card.Name,
		DomainId:    card.DomainID,
		Level:       int32(card.Level),
		Type:        domainCardTypeToProto(card.Type),
		RecallCost:  int32(card.RecallCost),
		UsageLimit:  card.UsageLimit,
		FeatureText: card.FeatureText,
	}
}

func toProtoDaggerheartDomainCards(cards []storage.DaggerheartDomainCard) []*pb.DaggerheartDomainCard {
	items := make([]*pb.DaggerheartDomainCard, 0, len(cards))
	for _, card := range cards {
		items = append(items, toProtoDaggerheartDomainCard(card))
	}
	return items
}

func toProtoDaggerheartWeapon(weapon storage.DaggerheartWeapon) *pb.DaggerheartWeapon {
	return &pb.DaggerheartWeapon{
		Id:         weapon.ID,
		Name:       weapon.Name,
		Category:   weaponCategoryToProto(weapon.Category),
		Tier:       int32(weapon.Tier),
		Trait:      weapon.Trait,
		Range:      weapon.Range,
		DamageDice: toProtoDaggerheartDamageDice(weapon.DamageDice),
		DamageType: damageTypeToProto(weapon.DamageType),
		Burden:     int32(weapon.Burden),
		Feature:    weapon.Feature,
	}
}

func toProtoDaggerheartWeapons(weapons []storage.DaggerheartWeapon) []*pb.DaggerheartWeapon {
	items := make([]*pb.DaggerheartWeapon, 0, len(weapons))
	for _, weapon := range weapons {
		items = append(items, toProtoDaggerheartWeapon(weapon))
	}
	return items
}

func toProtoDaggerheartArmor(armor storage.DaggerheartArmor) *pb.DaggerheartArmor {
	return &pb.DaggerheartArmor{
		Id:                  armor.ID,
		Name:                armor.Name,
		Tier:                int32(armor.Tier),
		BaseMajorThreshold:  int32(armor.BaseMajorThreshold),
		BaseSevereThreshold: int32(armor.BaseSevereThreshold),
		ArmorScore:          int32(armor.ArmorScore),
		Feature:             armor.Feature,
	}
}

func toProtoDaggerheartArmorList(items []storage.DaggerheartArmor) []*pb.DaggerheartArmor {
	armor := make([]*pb.DaggerheartArmor, 0, len(items))
	for _, item := range items {
		armor = append(armor, toProtoDaggerheartArmor(item))
	}
	return armor
}

func toProtoDaggerheartItem(item storage.DaggerheartItem) *pb.DaggerheartItem {
	return &pb.DaggerheartItem{
		Id:          item.ID,
		Name:        item.Name,
		Rarity:      itemRarityToProto(item.Rarity),
		Kind:        itemKindToProto(item.Kind),
		StackMax:    int32(item.StackMax),
		Description: item.Description,
		EffectText:  item.EffectText,
	}
}

func toProtoDaggerheartItems(items []storage.DaggerheartItem) []*pb.DaggerheartItem {
	results := make([]*pb.DaggerheartItem, 0, len(items))
	for _, item := range items {
		results = append(results, toProtoDaggerheartItem(item))
	}
	return results
}

func toProtoDaggerheartEnvironment(env storage.DaggerheartEnvironment) *pb.DaggerheartEnvironment {
	return &pb.DaggerheartEnvironment{
		Id:                    env.ID,
		Name:                  env.Name,
		Tier:                  int32(env.Tier),
		Type:                  environmentTypeToProto(env.Type),
		Difficulty:            int32(env.Difficulty),
		Impulses:              append([]string{}, env.Impulses...),
		PotentialAdversaryIds: append([]string{}, env.PotentialAdversaryIDs...),
		Features:              toProtoDaggerheartFeatures(env.Features),
		Prompts:               append([]string{}, env.Prompts...),
	}
}

func toProtoDaggerheartEnvironments(envs []storage.DaggerheartEnvironment) []*pb.DaggerheartEnvironment {
	results := make([]*pb.DaggerheartEnvironment, 0, len(envs))
	for _, env := range envs {
		results = append(results, toProtoDaggerheartEnvironment(env))
	}
	return results
}

func toProtoDaggerheartFeatures(features []storage.DaggerheartFeature) []*pb.DaggerheartFeature {
	items := make([]*pb.DaggerheartFeature, 0, len(features))
	for _, feature := range features {
		items = append(items, &pb.DaggerheartFeature{
			Id:          feature.ID,
			Name:        feature.Name,
			Description: feature.Description,
			Level:       int32(feature.Level),
		})
	}
	return items
}

func toProtoDaggerheartHopeFeature(feature storage.DaggerheartHopeFeature) *pb.DaggerheartHopeFeature {
	return &pb.DaggerheartHopeFeature{
		Name:        feature.Name,
		Description: feature.Description,
		HopeCost:    int32(feature.HopeCost),
	}
}

func toProtoDaggerheartDamageDice(dice []storage.DaggerheartDamageDie) []*pb.DiceSpec {
	results := make([]*pb.DiceSpec, 0, len(dice))
	for _, die := range dice {
		results = append(results, &pb.DiceSpec{
			Sides: int32(die.Sides),
			Count: int32(die.Count),
		})
	}
	return results
}

func heritageKindToProto(kind string) pb.DaggerheartHeritageKind {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "ancestry":
		return pb.DaggerheartHeritageKind_DAGGERHEART_HERITAGE_KIND_ANCESTRY
	case "community":
		return pb.DaggerheartHeritageKind_DAGGERHEART_HERITAGE_KIND_COMMUNITY
	default:
		return pb.DaggerheartHeritageKind_DAGGERHEART_HERITAGE_KIND_UNSPECIFIED
	}
}

func domainCardTypeToProto(kind string) pb.DaggerheartDomainCardType {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "ability":
		return pb.DaggerheartDomainCardType_DAGGERHEART_DOMAIN_CARD_TYPE_ABILITY
	case "spell":
		return pb.DaggerheartDomainCardType_DAGGERHEART_DOMAIN_CARD_TYPE_SPELL
	case "grimoire":
		return pb.DaggerheartDomainCardType_DAGGERHEART_DOMAIN_CARD_TYPE_GRIMOIRE
	default:
		return pb.DaggerheartDomainCardType_DAGGERHEART_DOMAIN_CARD_TYPE_UNSPECIFIED
	}
}

func weaponCategoryToProto(kind string) pb.DaggerheartWeaponCategory {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "primary":
		return pb.DaggerheartWeaponCategory_DAGGERHEART_WEAPON_CATEGORY_PRIMARY
	case "secondary":
		return pb.DaggerheartWeaponCategory_DAGGERHEART_WEAPON_CATEGORY_SECONDARY
	default:
		return pb.DaggerheartWeaponCategory_DAGGERHEART_WEAPON_CATEGORY_UNSPECIFIED
	}
}

func itemRarityToProto(kind string) pb.DaggerheartItemRarity {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "common":
		return pb.DaggerheartItemRarity_DAGGERHEART_ITEM_RARITY_COMMON
	case "uncommon":
		return pb.DaggerheartItemRarity_DAGGERHEART_ITEM_RARITY_UNCOMMON
	case "rare":
		return pb.DaggerheartItemRarity_DAGGERHEART_ITEM_RARITY_RARE
	case "unique":
		return pb.DaggerheartItemRarity_DAGGERHEART_ITEM_RARITY_UNIQUE
	default:
		return pb.DaggerheartItemRarity_DAGGERHEART_ITEM_RARITY_UNSPECIFIED
	}
}

func itemKindToProto(kind string) pb.DaggerheartItemKind {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "consumable":
		return pb.DaggerheartItemKind_DAGGERHEART_ITEM_KIND_CONSUMABLE
	case "equipment":
		return pb.DaggerheartItemKind_DAGGERHEART_ITEM_KIND_EQUIPMENT
	case "treasure":
		return pb.DaggerheartItemKind_DAGGERHEART_ITEM_KIND_TREASURE
	default:
		return pb.DaggerheartItemKind_DAGGERHEART_ITEM_KIND_UNSPECIFIED
	}
}

func environmentTypeToProto(kind string) pb.DaggerheartEnvironmentType {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "exploration":
		return pb.DaggerheartEnvironmentType_DAGGERHEART_ENVIRONMENT_TYPE_EXPLORATION
	case "social":
		return pb.DaggerheartEnvironmentType_DAGGERHEART_ENVIRONMENT_TYPE_SOCIAL
	case "traversal":
		return pb.DaggerheartEnvironmentType_DAGGERHEART_ENVIRONMENT_TYPE_TRAVERSAL
	case "event":
		return pb.DaggerheartEnvironmentType_DAGGERHEART_ENVIRONMENT_TYPE_EVENT
	default:
		return pb.DaggerheartEnvironmentType_DAGGERHEART_ENVIRONMENT_TYPE_UNSPECIFIED
	}
}

func damageTypeToProto(kind string) pb.DaggerheartDamageType {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "physical":
		return pb.DaggerheartDamageType_DAGGERHEART_DAMAGE_TYPE_PHYSICAL
	case "magic":
		return pb.DaggerheartDamageType_DAGGERHEART_DAMAGE_TYPE_MAGIC
	case "mixed":
		return pb.DaggerheartDamageType_DAGGERHEART_DAMAGE_TYPE_MIXED
	default:
		return pb.DaggerheartDamageType_DAGGERHEART_DAMAGE_TYPE_UNSPECIFIED
	}
}
