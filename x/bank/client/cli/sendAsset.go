package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/bank/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//SendAssetCmd : create a send asset tx and sign it with the give key
func SendAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendAsset",
		Short: "Sends an asset peg to an order transaction with a given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := context2.NewTxContextFromCLI().
				WithCodec(cdc)

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)

			msg := client.BuildSendAssetMsg(from, to, pegHashHex)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	return cmd
}
