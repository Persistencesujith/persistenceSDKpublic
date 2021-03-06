package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/utils"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/client/context"
	sdk "github.com/commitHub/commitBlockchain/types"

	"github.com/commitHub/commitBlockchain/wire"

	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	"github.com/commitHub/commitBlockchain/x/ibc"
)

//IBCSendAssetCmd : create a send asset tx and sign it with the give key
func IBCSendAssetCmd(cdc *wire.Codec) *cobra.Command {
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

			sourceChain := viper.GetString(client.FlagChainID)
			destinationChain := viper.GetString(flagChain)

			pegHashStr := viper.GetString(flagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)

			msg := ibc.BuildSendAssetMsg(from, to, pegHashHex, sourceChain, destinationChain)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagTo, "", "Address to issue fiat to")
	cmd.Flags().String(flagChain, "", "Destination chain to send coins")
	cmd.Flags().String(flagPegHash, "", "Peg Hash")
	return cmd
}
