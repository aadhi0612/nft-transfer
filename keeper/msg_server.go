package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/aadhi0612/nft-transfer/types"
)

var _ types.MsgServer = Keeper{}

// Transfer defines a rpc handler method for MsgTransfer.
func (k Keeper) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	sequence, err := k.SendTransfer(
		ctx, msg.SourcePort, msg.SourceChannel, msg.ClassId, msg.TokenIds,
		sender, msg.Receiver, msg.TimeoutHeight, msg.TimeoutTimestamp, msg.Memo,
	)
	if err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("IBC non-fungible token transfer",
		"classID", msg.ClassId,
		"tokenIDs", strings.Join(msg.TokenIds, ","),
		"sender", msg.Sender,
		"receiver", msg.Receiver,
	)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgTransferResponse{Sequence: sequence}, nil
}

// UpdateParams defines a governance operation for updating the nft-transfer module parameters.
// The authority is defined in the keeper.
func (k Keeper) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.GetAuthority() != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}
