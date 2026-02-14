package scenario

import (
	"context"

	gamev1 "github.com/louisbranch/fracturing.space/api/gen/go/game/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unimplemented(method string) error {
	return status.Errorf(codes.Unimplemented, "%s not implemented", method)
}

type fakeCampaignClient struct {
	create func(context.Context, *gamev1.CreateCampaignRequest, ...grpc.CallOption) (*gamev1.CreateCampaignResponse, error)
}

func (f *fakeCampaignClient) CreateCampaign(ctx context.Context, in *gamev1.CreateCampaignRequest, opts ...grpc.CallOption) (*gamev1.CreateCampaignResponse, error) {
	if f.create != nil {
		return f.create(ctx, in, opts...)
	}
	return nil, unimplemented("CreateCampaign")
}

func (f *fakeCampaignClient) ListCampaigns(context.Context, *gamev1.ListCampaignsRequest, ...grpc.CallOption) (*gamev1.ListCampaignsResponse, error) {
	return nil, unimplemented("ListCampaigns")
}

func (f *fakeCampaignClient) GetCampaign(context.Context, *gamev1.GetCampaignRequest, ...grpc.CallOption) (*gamev1.GetCampaignResponse, error) {
	return nil, unimplemented("GetCampaign")
}

func (f *fakeCampaignClient) EndCampaign(context.Context, *gamev1.EndCampaignRequest, ...grpc.CallOption) (*gamev1.EndCampaignResponse, error) {
	return nil, unimplemented("EndCampaign")
}

func (f *fakeCampaignClient) ArchiveCampaign(context.Context, *gamev1.ArchiveCampaignRequest, ...grpc.CallOption) (*gamev1.ArchiveCampaignResponse, error) {
	return nil, unimplemented("ArchiveCampaign")
}

func (f *fakeCampaignClient) RestoreCampaign(context.Context, *gamev1.RestoreCampaignRequest, ...grpc.CallOption) (*gamev1.RestoreCampaignResponse, error) {
	return nil, unimplemented("RestoreCampaign")
}

type fakeParticipantClient struct {
	create func(context.Context, *gamev1.CreateParticipantRequest, ...grpc.CallOption) (*gamev1.CreateParticipantResponse, error)
}

func (f *fakeParticipantClient) CreateParticipant(ctx context.Context, in *gamev1.CreateParticipantRequest, opts ...grpc.CallOption) (*gamev1.CreateParticipantResponse, error) {
	if f.create != nil {
		return f.create(ctx, in, opts...)
	}
	return nil, unimplemented("CreateParticipant")
}

func (f *fakeParticipantClient) UpdateParticipant(context.Context, *gamev1.UpdateParticipantRequest, ...grpc.CallOption) (*gamev1.UpdateParticipantResponse, error) {
	return nil, unimplemented("UpdateParticipant")
}

func (f *fakeParticipantClient) DeleteParticipant(context.Context, *gamev1.DeleteParticipantRequest, ...grpc.CallOption) (*gamev1.DeleteParticipantResponse, error) {
	return nil, unimplemented("DeleteParticipant")
}

func (f *fakeParticipantClient) ListParticipants(context.Context, *gamev1.ListParticipantsRequest, ...grpc.CallOption) (*gamev1.ListParticipantsResponse, error) {
	return nil, unimplemented("ListParticipants")
}

func (f *fakeParticipantClient) GetParticipant(context.Context, *gamev1.GetParticipantRequest, ...grpc.CallOption) (*gamev1.GetParticipantResponse, error) {
	return nil, unimplemented("GetParticipant")
}

type fakeCharacterClient struct {
	create            func(context.Context, *gamev1.CreateCharacterRequest, ...grpc.CallOption) (*gamev1.CreateCharacterResponse, error)
	setDefaultControl func(context.Context, *gamev1.SetDefaultControlRequest, ...grpc.CallOption) (*gamev1.SetDefaultControlResponse, error)
	patchProfile      func(context.Context, *gamev1.PatchCharacterProfileRequest, ...grpc.CallOption) (*gamev1.PatchCharacterProfileResponse, error)
	patchState        func(context.Context, *gamev1.PatchCharacterStateRequest, ...grpc.CallOption) (*gamev1.PatchCharacterStateResponse, error)
}

func (f *fakeCharacterClient) CreateCharacter(ctx context.Context, in *gamev1.CreateCharacterRequest, opts ...grpc.CallOption) (*gamev1.CreateCharacterResponse, error) {
	if f.create != nil {
		return f.create(ctx, in, opts...)
	}
	return nil, unimplemented("CreateCharacter")
}

func (f *fakeCharacterClient) UpdateCharacter(context.Context, *gamev1.UpdateCharacterRequest, ...grpc.CallOption) (*gamev1.UpdateCharacterResponse, error) {
	return nil, unimplemented("UpdateCharacter")
}

func (f *fakeCharacterClient) DeleteCharacter(context.Context, *gamev1.DeleteCharacterRequest, ...grpc.CallOption) (*gamev1.DeleteCharacterResponse, error) {
	return nil, unimplemented("DeleteCharacter")
}

func (f *fakeCharacterClient) ListCharacters(context.Context, *gamev1.ListCharactersRequest, ...grpc.CallOption) (*gamev1.ListCharactersResponse, error) {
	return nil, unimplemented("ListCharacters")
}

func (f *fakeCharacterClient) SetDefaultControl(ctx context.Context, in *gamev1.SetDefaultControlRequest, opts ...grpc.CallOption) (*gamev1.SetDefaultControlResponse, error) {
	if f.setDefaultControl != nil {
		return f.setDefaultControl(ctx, in, opts...)
	}
	return nil, unimplemented("SetDefaultControl")
}

func (f *fakeCharacterClient) GetCharacterSheet(context.Context, *gamev1.GetCharacterSheetRequest, ...grpc.CallOption) (*gamev1.GetCharacterSheetResponse, error) {
	return nil, unimplemented("GetCharacterSheet")
}

func (f *fakeCharacterClient) PatchCharacterProfile(ctx context.Context, in *gamev1.PatchCharacterProfileRequest, opts ...grpc.CallOption) (*gamev1.PatchCharacterProfileResponse, error) {
	if f.patchProfile != nil {
		return f.patchProfile(ctx, in, opts...)
	}
	return nil, unimplemented("PatchCharacterProfile")
}

func (f *fakeCharacterClient) PatchCharacterState(ctx context.Context, in *gamev1.PatchCharacterStateRequest, opts ...grpc.CallOption) (*gamev1.PatchCharacterStateResponse, error) {
	if f.patchState != nil {
		return f.patchState(ctx, in, opts...)
	}
	return nil, unimplemented("PatchCharacterState")
}

type fakeEventClient struct {
	seq int64
}

func (f *fakeEventClient) AppendEvent(context.Context, *gamev1.AppendEventRequest, ...grpc.CallOption) (*gamev1.AppendEventResponse, error) {
	return nil, unimplemented("AppendEvent")
}

func (f *fakeEventClient) ListEvents(context.Context, *gamev1.ListEventsRequest, ...grpc.CallOption) (*gamev1.ListEventsResponse, error) {
	f.seq++
	return &gamev1.ListEventsResponse{
		Events: []*gamev1.Event{{Seq: uint64(f.seq)}},
	}, nil
}

type fakeSnapshotClient struct {
	patchState func(context.Context, *gamev1.PatchCharacterStateRequest, ...grpc.CallOption) (*gamev1.PatchCharacterStateResponse, error)
}

func (f *fakeSnapshotClient) GetSnapshot(context.Context, *gamev1.GetSnapshotRequest, ...grpc.CallOption) (*gamev1.GetSnapshotResponse, error) {
	return nil, unimplemented("GetSnapshot")
}

func (f *fakeSnapshotClient) PatchCharacterState(ctx context.Context, in *gamev1.PatchCharacterStateRequest, opts ...grpc.CallOption) (*gamev1.PatchCharacterStateResponse, error) {
	if f.patchState != nil {
		return f.patchState(ctx, in, opts...)
	}
	return nil, unimplemented("PatchCharacterState")
}

func (f *fakeSnapshotClient) UpdateSnapshotState(context.Context, *gamev1.UpdateSnapshotStateRequest, ...grpc.CallOption) (*gamev1.UpdateSnapshotStateResponse, error) {
	return nil, unimplemented("UpdateSnapshotState")
}
