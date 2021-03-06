package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"
	cliclient "github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	"github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/rest"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	authctx "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/fiatFactory"
)

//ExecuteFiatHandlerFunction : handles execute fiat rest message
func ExecuteFiatHandlerFunction(cliCtx context.CLIContext, cdc *wire.Codec, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var msg fiatFactory.ExecuteFiatBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = msgWireCdc.UnmarshalJSON(body, &msg)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = govalidator.ValidateStruct(msg)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		adjustment, ok := utils.ParseFloat64OrReturnBadRequest(w, msg.GasAdjustment, cliclient.DefaultGasAdjustment)
		if !ok {
			return
		}

		cliCtx = cliCtx.WithGasAdjustment(adjustment)
		cliCtx = cliCtx.WithFromAddressName(msg.From)
		cliCtx.JSON = true

		if err := cliCtx.EnsureAccountExists(); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		from, err := cliCtx.GetFromAddress()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ownerAddressStr := msg.OwnerAddress

		ownerAddress, err := sdk.AccAddressFromBech32(ownerAddressStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		toStr := msg.To

		to, err := sdk.AccAddressFromBech32(toStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		txCtx := authctx.TxContext{
			Codec:         cdc,
			AccountNumber: msg.AccountNumber,
			Sequence:      msg.Sequence,
			Gas:           msg.Gas,
			ChainID:       msg.ChainID,
		}
		assetPegHashHex, err := sdk.GetAssetPegHashHex(msg.AssetPegHash)
		fiatPegHashHex, err := sdk.GetFiatPegHashHex(msg.FiatPegHash)
		fiatPeg := sdk.BaseFiatPeg{
			PegHash:           fiatPegHashHex,
			TransactionAmount: msg.Amount,
		}

		buildmsg := fiatFactory.BuildExecuteFiatMsg(from, ownerAddress, to, assetPegHashHex, sdk.FiatPegWallet{fiatPeg})

		if kafka == true {
			ticketID := rest.TicketIDGenerator("FFEF")

			jsonResponse := rest.SendToKafka(rest.NewKafkaMsgFromRest(buildmsg, ticketID, txCtx, cliCtx, msg.Password), kafkaState, cdc)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{buildmsg}, msg.Password)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			w.Write(utils.ResponseBytesToJSON(output))
		}
	}

}
