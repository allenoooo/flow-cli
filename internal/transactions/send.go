/*
 * Flow CLI
 *
 * Copyright 2019-2021 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package transactions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
)

type flagsSend struct {
	ArgsJSON string   `default:"" flag:"args-json" info:"arguments in JSON-Cadence format"`
	Arg      []string `default:"" flag:"arg" info:"argument in Type:Value format"`
	Signer   string   `default:"emulator-account" flag:"signer" info:"Account name from configuration used to sign the transaction"`
	GasLimit uint64   `default:"1000" flag:"gas-limit" info:"transaction gas limit"`
	Include  []string `default:"" flag:"include" info:"Fields to include in the output"`
	Exclude  []string `default:"" flag:"exclude" info:"Fields to exclude from the output (events)"`
}

var sendFlags = flagsSend{}

var SendCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:     "send <code filename>",
		Short:   "Send a transaction",
		Args:    cobra.ExactArgs(1),
		Example: `flow transactions send tx.cdc --arg String:"Hello world"`,
	},
	Flags: &sendFlags,
	RunS:  send,
}

func send(
	args []string,
	readerWriter flowkit.ReaderWriter,
	globalFlags command.GlobalFlags,
	services *services.Services,
	state *flowkit.State,
) (command.Result, error) {
	codeFilename := args[0]

	signer, err := state.Accounts().ByName(sendFlags.Signer)
	if err != nil {
		return nil, err
	}

	code, err := readerWriter.ReadFile(codeFilename)
	if err != nil {
		return nil, fmt.Errorf("error loading transaction file: %w", err)
	}

	txArgs, err := flowkit.ParseArguments(sendFlags.Arg, sendFlags.ArgsJSON)
	if err != nil {
		return nil, err
	}

	tx, result, err := services.Transactions.Send(
		signer,
		code,
		codeFilename,
		sendFlags.GasLimit,
		txArgs,
		globalFlags.Network,
	)
	if err != nil {
		return nil, err
	}

	return &TransactionResult{
		result:  result,
		tx:      tx,
		include: sendFlags.Include,
		exclude: sendFlags.Exclude,
	}, nil
}
