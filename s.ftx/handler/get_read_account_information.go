package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	"swallowtail/s.ftx/marshaling"
	ftxproto "swallowtail/s.ftx/proto"
)

// ReadAccountInformation ...
func (S *FTXService) ReadAccountInformation(
	ctx context.Context, in *ftxproto.ReadAccountInformationRequest,
) (*ftxproto.ReadAccountInformationResponse, error) {
	// Basic validation.
	switch {
	case in.GetCredentials() == nil:
		return nil, gerrors.BadParam("missing_param.credentials", nil)
	}

	// Validate credentials.
	if err := validateCredentials(in.Credentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_account_information", nil)
	}

	errParams := map[string]string{
		"subaccount": in.GetCredentials().Subaccount,
		"url":        in.GetCredentials().Url,
		"ws_url":     in.GetCredentials().WsUrl,
	}

	// Marshal to domain.
	domainCredentials := marshaling.VenueCredentialsProtoToFTXCredentials(in.Credentials)

	// Read account information.
	rsp, err := client.ReadAccountInformation(ctx, &client.ReadAccountInformationRequest{}, domainCredentials)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_account_information", errParams)
	}

	// Validate response.
	if !rsp.Success {
		return nil, gerrors.Augment(err, "failed_to_read_account_information.ftx_client_failure", errParams)
	}
	if rsp.Result == nil {
		return nil, gerrors.Augment(err, "failed_to_read_account_information.ftx_client_nil_result", errParams)
	}

	// Marshal to proto.
	return marshaling.ReadAccountInformationDomainToProto(rsp.Result), nil
}
