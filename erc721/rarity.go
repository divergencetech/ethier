// Package rarity implements computation of ERC721 metadata rarity based on
// information entropy.
//
// Every trait/attribute of an ERC721 NFT can be considered as a categorical
// probability distribution. Each possible
package erc721

import (
	"math"
)

// Rarity computes rarity of each token in the Collection based on information
// entropy. Every TraitType is considered as a categorical probability
// distribution with each Value having an associated probability and hence
// information content. The rarity of a particular token is the sum of
// information content carried by each of its Attributes.
//
// Notably, the lack of a TraitType is considered as a null-Value Attribute as
// the absence across the majority of a Collection implies rarity in those
// tokens that do carry the TraitType.
//
// Non-string Attribute Values are passed to the bucket function. The returned
// bucket is used in place of original value. It is valid for the bucket
// function to simply return the string equivalent (e.g. true/false for
// booleans).
//
// The returned rarity scores are normalised to a maximum value of 1 across the
// Collection. Without this, inter-Collection comparisons would be invalid as
// those that simply carry more Attributes per token would have higher scores.
func (coll Collection) Rarity(bucket func(interface{}) string) map[TokenID]float64 {

	// distribution and counts carry floats instead of integers to make
	// calculation of entropy simpler. counts[x] will contain the sum of all
	// values in distribution[x], which is split by value.
	distributions := make(map[string]map[string]float64)
	counts := make(map[string]float64)
	attributes := make(map[TokenID]map[string]string)

	for id, meta := range coll {
		attributes[id] = make(map[string]string)
		for _, attr := range meta.Attributes {
			counts[attr.TraitType]++

			val, ok := attr.Value.(string)
			if !ok {
				// We can't make any assumptions about non-string discrete types
				// (e.g. booleans) because they can technically clash with
				// equivalent string values. Lesson: strong typing is important
				// and JSON is stupid.
				val = bucket(attr.Value)
			}
			attributes[id][attr.TraitType] = val

			if _, ok := distributions[attr.TraitType]; !ok {
				distributions[attr.TraitType] = make(map[string]float64)
			}
			distributions[attr.TraitType][val]++
		}
	}

	var max float64
	scores := make(map[TokenID]float64)
	collSize := float64(len(coll))
	for id := range coll {
		// It's important to calculate over all possible attributes, even those
		// that a particular token lacks. Without this, we would favour tokens
		// that simply have more traits.
		for attr := range counts {
			var n float64
			if v, ok := attributes[id][attr]; ok {
				n = distributions[attr][v]
			} else {
				n = collSize - counts[attr]
			}

			scores[id] += -math.Log2(n / collSize)
		}

		max = math.Max(max, scores[id])
	}

	for id := range scores {
		scores[id] /= max
	}
	return scores
}
