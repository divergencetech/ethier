package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/divergencetech/ethier/erc721"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "rarity",
		Short: `Calculates information-theoretic ERC721 metadata rarity, reading a JSON list from stdin, and printing ranked rarity scores on stdout.`,
		RunE:  rarity,
	}

	rootCmd.AddCommand(cmd)
}

func rarity(cmd *cobra.Command, args []string) error {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("io.ReadAll(stdin): %v", err)
	}

	var md []*erc721.Metadata
	if err := json.Unmarshal(buf, &md); err != nil {
		return fmt.Errorf("json.Unmarshal(stdin, %T): %v", md, err)
	}

	type pair struct {
		id    erc721.TokenID
		score float64
	}
	var scores []pair
	res := erc721.CollectionFromMetadata(md).Rarity(func(x interface{}) string {
		log.Fatalf("Non-string attribute of type %T not supported by CLI", x)
		return ""
	})

	for id, score := range res.Scores {
		scores = append(scores, pair{id, score})
	}
	sort.Slice(scores, func(i, j int) bool {
		sI, sJ := scores[i], scores[j]
		if scI, scJ := sI.score, sJ.score; scI != scJ {
			return scI > scJ
		}
		return sI.id.Cmp(&sJ.id) == -1
	})

	fmt.Fprintf(os.Stderr, "Collection entropy: %.4f\n", res.Entropy)
	for _, s := range scores {
		fmt.Printf("%s %.4f\n", s.id.String(), s.score)
	}
	return nil
}
