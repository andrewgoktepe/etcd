// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"fmt"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	getConsistency string
	getLimit       int64
	getSortOrder   string
	getSortTarget  string
	getPrefix      bool
	getFromKey     bool
)

// NewGetCommand returns the cobra command for "get".
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [options] <key> [range_end]",
		Short: "Get gets the key or a range of keys.",
		Run:   getCommandFunc,
	}

	cmd.Flags().StringVar(&getConsistency, "consistency", "l", "Linearizable(l) or Serializable(s)")
	cmd.Flags().StringVar(&getSortOrder, "order", "", "order of results; ASCEND or DESCEND")
	cmd.Flags().StringVar(&getSortTarget, "sort-by", "", "sort target; CREATE, KEY, MODIFY, VALUE, or VERSION")
	cmd.Flags().Int64Var(&getLimit, "limit", 0, "maximum number of results")
	cmd.Flags().BoolVar(&getPrefix, "prefix", false, "get keys with matching prefix")
	cmd.Flags().BoolVar(&getFromKey, "from-key", false, "get keys that are greater than or equal to the given key")
	return cmd
}

// getCommandFunc executes the "get" command.
func getCommandFunc(cmd *cobra.Command, args []string) {
	key, opts := getGetOp(cmd, args)
	resp, err := mustClientFromCmd(cmd).Get(context.TODO(), key, opts...)
	if err != nil {
		ExitWithError(ExitError, err)
	}

	display.Get(*resp)
}

func getGetOp(cmd *cobra.Command, args []string) (string, []clientv3.OpOption) {
	if len(args) == 0 {
		ExitWithError(ExitBadArgs, fmt.Errorf("range command needs arguments."))
	}

	if getPrefix && getFromKey {
		ExitWithError(ExitBadArgs, fmt.Errorf("`--prefix` and `--from-key` cannot be set at the same time, choose one."))
	}

	opts := []clientv3.OpOption{}
	switch getConsistency {
	case "s":
		opts = append(opts, clientv3.WithSerializable())
	case "l":
	default:
		ExitWithError(ExitBadFeature, fmt.Errorf("unknown consistency flag %q", getConsistency))
	}

	key := args[0]
	if len(args) > 1 {
		if getPrefix || getFromKey {
			ExitWithError(ExitBadArgs, fmt.Errorf("too many arguments, only accept one arguement when `--prefix` or `--from-key` is set."))
		}
		opts = append(opts, clientv3.WithRange(args[1]))
	}

	opts = append(opts, clientv3.WithLimit(getLimit))

	sortByOrder := clientv3.SortNone
	sortOrder := strings.ToUpper(getSortOrder)
	switch {
	case sortOrder == "ASCEND":
		sortByOrder = clientv3.SortAscend
	case sortOrder == "DESCEND":
		sortByOrder = clientv3.SortDescend
	case sortOrder == "":
		// nothing
	default:
		ExitWithError(ExitBadFeature, fmt.Errorf("bad sort order %v", getSortOrder))
	}

	sortByTarget := clientv3.SortByKey
	sortTarget := strings.ToUpper(getSortTarget)
	switch {
	case sortTarget == "CREATE":
		sortByTarget = clientv3.SortByCreatedRev
	case sortTarget == "KEY":
		sortByTarget = clientv3.SortByKey
	case sortTarget == "MODIFY":
		sortByTarget = clientv3.SortByModRevision
	case sortTarget == "VALUE":
		sortByTarget = clientv3.SortByValue
	case sortTarget == "VERSION":
		sortByTarget = clientv3.SortByVersion
	case sortTarget == "":
		// nothing
	default:
		ExitWithError(ExitBadFeature, fmt.Errorf("bad sort target %v", getSortTarget))
	}

	opts = append(opts, clientv3.WithSort(sortByTarget, sortByOrder))

	if getPrefix {
		opts = append(opts, clientv3.WithPrefix())
	}

	if getFromKey {
		opts = append(opts, clientv3.WithFromKey())
	}

	return key, opts
}
