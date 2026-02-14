package game

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	campaignv1 "github.com/louisbranch/fracturing.space/api/gen/go/game/v1"
	grpcmeta "github.com/louisbranch/fracturing.space/internal/services/game/api/grpc/metadata"
	"github.com/louisbranch/fracturing.space/internal/services/game/domain/campaign"
	"github.com/louisbranch/fracturing.space/internal/services/game/domain/campaign/character"
	"github.com/louisbranch/fracturing.space/internal/services/game/domain/campaign/event"
	"github.com/louisbranch/fracturing.space/internal/services/game/domain/systems/daggerheart"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type characterCreator struct {
	stores      Stores
	clock       func() time.Time
	idGenerator func() (string, error)
}

func newCharacterCreator(service *CharacterService) characterCreator {
	creator := characterCreator{stores: service.stores, clock: service.clock, idGenerator: service.idGenerator}
	if creator.clock == nil {
		creator.clock = time.Now
	}
	return creator
}

func (c characterCreator) create(ctx context.Context, campaignID string, in *campaignv1.CreateCharacterRequest) (character.Character, error) {
	campaignRecord, err := c.stores.Campaign.Get(ctx, campaignID)
	if err != nil {
		return character.Character{}, err
	}
	if err := campaign.ValidateCampaignOperation(campaignRecord.Status, campaign.CampaignOpRead); err != nil {
		return character.Character{}, err
	}

	input := character.CreateCharacterInput{
		CampaignID: campaignID,
		Name:       in.GetName(),
		Kind:       characterKindFromProto(in.GetKind()),
		Notes:      in.GetNotes(),
	}
	normalized, err := character.NormalizeCreateCharacterInput(input)
	if err != nil {
		return character.Character{}, err
	}

	characterID, err := c.idGenerator()
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "generate character id: %v", err)
	}

	payload := event.CharacterCreatedPayload{
		CharacterID: characterID,
		Name:        normalized.Name,
		Kind:        in.GetKind().String(),
		Notes:       normalized.Notes,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "encode payload: %v", err)
	}

	actorID := grpcmeta.ParticipantIDFromContext(ctx)
	actorType := event.ActorTypeSystem
	if actorID != "" {
		actorType = event.ActorTypeParticipant
	}

	stored, err := c.stores.Event.AppendEvent(ctx, event.Event{
		CampaignID:   campaignID,
		Timestamp:    c.clock().UTC(),
		Type:         event.TypeCharacterCreated,
		RequestID:    grpcmeta.RequestIDFromContext(ctx),
		InvocationID: grpcmeta.InvocationIDFromContext(ctx),
		ActorType:    actorType,
		ActorID:      actorID,
		EntityType:   "character",
		EntityID:     characterID,
		PayloadJSON:  payloadJSON,
	})
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "append event: %v", err)
	}

	applier := c.stores.Applier()
	if err := applier.Apply(ctx, stored); err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "apply event: %v", err)
	}

	created, err := c.stores.Character.GetCharacter(ctx, campaignID, characterID)
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "load character: %v", err)
	}

	// Get Daggerheart defaults for the character kind
	kindStr := "PC"
	if created.Kind == character.CharacterKindNPC {
		kindStr = "NPC"
	}
	dhDefaults := daggerheart.GetProfileDefaults(kindStr)

	reqID := grpcmeta.RequestIDFromContext(ctx)
	invocationID := grpcmeta.InvocationIDFromContext(ctx)
	profileActorType := event.ActorTypeSystem
	if actorID != "" {
		profileActorType = event.ActorTypeGM
	}

	experiencesPayload := make([]map[string]any, 0, len(dhDefaults.Experiences))
	for _, experience := range dhDefaults.Experiences {
		experiencesPayload = append(experiencesPayload, map[string]any{
			"name":     experience.Name,
			"modifier": experience.Modifier,
		})
	}

	profilePayload := event.ProfileUpdatedPayload{
		CharacterID: created.ID,
		SystemProfile: map[string]any{
			"daggerheart": map[string]any{
				"level":            dhDefaults.Level,
				"hp_max":           dhDefaults.HpMax,
				"stress_max":       dhDefaults.StressMax,
				"evasion":          dhDefaults.Evasion,
				"major_threshold":  dhDefaults.MajorThreshold,
				"severe_threshold": dhDefaults.SevereThreshold,
				"proficiency":      dhDefaults.Proficiency,
				"armor_score":      dhDefaults.ArmorScore,
				"armor_max":        dhDefaults.ArmorMax,
				"agility":          dhDefaults.Traits.Agility,
				"strength":         dhDefaults.Traits.Strength,
				"finesse":          dhDefaults.Traits.Finesse,
				"instinct":         dhDefaults.Traits.Instinct,
				"presence":         dhDefaults.Traits.Presence,
				"knowledge":        dhDefaults.Traits.Knowledge,
				"experiences":      experiencesPayload,
			},
		},
	}
	profileJSON, err := json.Marshal(profilePayload)
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "encode profile payload: %v", err)
	}
	profileEvent, err := c.stores.Event.AppendEvent(ctx, event.Event{
		CampaignID:   campaignID,
		Timestamp:    c.clock().UTC(),
		Type:         event.TypeProfileUpdated,
		RequestID:    reqID,
		InvocationID: invocationID,
		ActorType:    profileActorType,
		ActorID:      actorID,
		EntityType:   "character",
		EntityID:     created.ID,
		PayloadJSON:  profileJSON,
	})
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "append profile event: %v", err)
	}

	hpAfter := dhDefaults.HpMax
	hopeAfter := daggerheart.HopeDefault
	hopeMaxAfter := daggerheart.HopeMax
	stressAfter := daggerheart.StressDefault
	armorAfter := 0
	lifeStateAfter := daggerheart.LifeStateAlive
	statePayload := daggerheart.CharacterStatePatchedPayload{
		CharacterID:    created.ID,
		HpAfter:        &hpAfter,
		HopeAfter:      &hopeAfter,
		HopeMaxAfter:   &hopeMaxAfter,
		StressAfter:    &stressAfter,
		ArmorAfter:     &armorAfter,
		LifeStateAfter: &lifeStateAfter,
	}
	stateJSON, err := json.Marshal(statePayload)
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "encode state payload: %v", err)
	}
	stateEvent, err := c.stores.Event.AppendEvent(ctx, event.Event{
		CampaignID:    campaignID,
		Timestamp:     c.clock().UTC(),
		Type:          daggerheart.EventTypeCharacterStatePatched,
		SessionID:     grpcmeta.SessionIDFromContext(ctx),
		RequestID:     reqID,
		InvocationID:  invocationID,
		ActorType:     profileActorType,
		ActorID:       actorID,
		EntityType:    "character",
		EntityID:      created.ID,
		SystemID:      campaignRecord.System.String(),
		SystemVersion: daggerheart.SystemVersion,
		PayloadJSON:   stateJSON,
	})
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "append state event: %v", err)
	}

	projectionApplier := c.stores.Applier()
	if err := projectionApplier.Apply(ctx, profileEvent); err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "apply profile event: %v", err)
	}
	adapter := daggerheart.NewAdapter(c.stores.Daggerheart)
	if err := adapter.ApplyEvent(ctx, stateEvent); err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "apply state event: %v", err)
	}

	return created, nil
}

type characterUpdater struct {
	stores Stores
	clock  func() time.Time
}

func newCharacterUpdater(service *CharacterService) characterUpdater {
	updater := characterUpdater{stores: service.stores, clock: service.clock}
	if updater.clock == nil {
		updater.clock = time.Now
	}
	return updater
}

func (u characterUpdater) update(ctx context.Context, campaignID string, in *campaignv1.UpdateCharacterRequest) (character.Character, error) {
	campaignRecord, err := u.stores.Campaign.Get(ctx, campaignID)
	if err != nil {
		return character.Character{}, err
	}
	if err := campaign.ValidateCampaignOperation(campaignRecord.Status, campaign.CampaignOpCampaignMutate); err != nil {
		return character.Character{}, err
	}

	characterID := strings.TrimSpace(in.GetCharacterId())
	if characterID == "" {
		return character.Character{}, status.Error(codes.InvalidArgument, "character id is required")
	}

	ch, err := u.stores.Character.GetCharacter(ctx, campaignID, characterID)
	if err != nil {
		return character.Character{}, err
	}

	fields := make(map[string]any)
	if name := in.GetName(); name != nil {
		trimmed := strings.TrimSpace(name.GetValue())
		if trimmed == "" {
			return character.Character{}, status.Error(codes.InvalidArgument, "name must not be empty")
		}
		ch.Name = trimmed
		fields["name"] = trimmed
	}
	if in.GetKind() != campaignv1.CharacterKind_CHARACTER_KIND_UNSPECIFIED {
		kind := characterKindFromProto(in.GetKind())
		if kind == character.CharacterKindUnspecified {
			return character.Character{}, status.Error(codes.InvalidArgument, "character kind is invalid")
		}
		ch.Kind = kind
		fields["kind"] = in.GetKind().String()
	}
	if notes := in.GetNotes(); notes != nil {
		ch.Notes = strings.TrimSpace(notes.GetValue())
		fields["notes"] = ch.Notes
	}
	if len(fields) == 0 {
		return character.Character{}, status.Error(codes.InvalidArgument, "at least one field must be provided")
	}

	payload := event.CharacterUpdatedPayload{
		CharacterID: characterID,
		Fields:      fields,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "encode payload: %v", err)
	}

	actorID := grpcmeta.ParticipantIDFromContext(ctx)
	actorType := event.ActorTypeSystem
	if actorID != "" {
		actorType = event.ActorTypeParticipant
	}

	stored, err := u.stores.Event.AppendEvent(ctx, event.Event{
		CampaignID:   campaignID,
		Timestamp:    u.clock().UTC(),
		Type:         event.TypeCharacterUpdated,
		RequestID:    grpcmeta.RequestIDFromContext(ctx),
		InvocationID: grpcmeta.InvocationIDFromContext(ctx),
		ActorType:    actorType,
		ActorID:      actorID,
		EntityType:   "character",
		EntityID:     characterID,
		PayloadJSON:  payloadJSON,
	})
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "append event: %v", err)
	}

	applier := u.stores.Applier()
	if err := applier.Apply(ctx, stored); err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "apply event: %v", err)
	}

	updated, err := u.stores.Character.GetCharacter(ctx, campaignID, characterID)
	if err != nil {
		return character.Character{}, status.Errorf(codes.Internal, "load character: %v", err)
	}

	return updated, nil
}
