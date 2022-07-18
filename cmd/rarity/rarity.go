// The rarity binary calculates ERC721 metadata rarity, reading a JSON list from
// stdin, and printing ranked rarity scores on stdout.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/divergencetech/ethier/erc721"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("io.ReadAll(stdin): %v", err)
	}

	var md []*erc721.Metadata
	if err := json.Unmarshal(buf, &md); err != nil {
		log.Fatalf("json.Unmarshal(stdin, %T): %v", md, err)
	}

	type pair struct {
		id    erc721.TokenID
		score float64
	}
	var scores []pair
	for id, score := range erc721.CollectionFromMetadata(md).Rarity(nil) {
		scores = append(scores, pair{id, score})
	}
	sort.Slice(scores, func(i, j int) bool {
		sI, sJ := scores[i], scores[j]
		if scI, scJ := sI.score, sJ.score; scI != scJ {
			return scI > scJ
		}
		return sI.id.Cmp(&sJ.id) == -1
	})

	for _, s := range scores {
		fmt.Println(s.id.String(), s.score)
	}
}
