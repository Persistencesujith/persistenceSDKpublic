package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/bank"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//DefineOrganizationCmd : Add the zone ID
func DefineOrganizationCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defineOrganization",
		Short: "define an organization address in acl",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := context2.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().WithAccountDecoder(authcmd.GetAccountDecoder(cdc)).WithCodec(cdc).WithLogger(os.Stdout)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			toStr := viper.GetString(FlagTo)
			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			strOrganizationID := viper.GetString(FlagOrganizationID)
			organizationID, err := sdk.GetOrganizationIDFromString(strOrganizationID)
			if err != nil {
				return nil
			}

			strZoneID := viper.GetString(FlagZoneID)
			zoneID, err := sdk.GetZoneIDFromString(strZoneID)
			if err != nil {
				return nil
			}

			msg := bank.BuildMsgDefineOrganization(from, to, organizationID, zoneID)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsOrganizationID)
	cmd.Flags().AddFlagSet(fsZoneID)
	cmd.Flags().AddFlagSet(fsTo)
	return cmd
}
