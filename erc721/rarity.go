package erc721

import (
	"math"
)

// Rarity describes the information-theoretic "rarity" of a Collection.
//
// The concept of "rarity" can be considered as a measure of "surprise" at the
// occurrence of a particular token's properties, within the context of the
// Collection from which it is derived. Self-information is a measure of such
// surprise, and information entropy a measure of the expected value of
// self-information across a distribution (i.e. across a Collection).
//
// It is trivial to "stuff" a Collection with extra information by merely adding
// additional properties to all tokens. This is reflected in the Entropy field,
// measured in bits—all else held equal, a Collection with more token properties
// will have higher Entropy. However, this information bloat is carried by the
// tokens themselves, so their individual information-content grows in line with
// Collection-wide Entropy. The Scores are therefore scaled down by the Entropy
// to provide unitless "relative surprise", which can be safely compared between
// Collections.
type Rarity struct {
	Entropy float64
	Scores  map[TokenID]float64
}

// Rarity computes rarity of each token in the Collection based on information
// entropy. Every TraitType is considered as a categorical probability
// distribution with each Value having an associated probability and hence
// information content. The rarity of a particular token is the sum of
// information content carried by each of its Attributes, divided by the entropy
// of the Collection as a whole (see the Rarity struct for rationale).
//
// Notably, the lack of a TraitType is considered as a null-Value Attribute as
// the absence across the majority of a Collection implies rarity in those
// tokens that do carry the TraitType.
//
// Non-string Attribute Values are passed to the bucket function. The returned
// bucket is used in place of original value. It is valid for the bucket
// function to simply return the string equivalent (e.g. true/false for
// booleans).
func (coll Collection) Rarity(bucket func(interface{}) string) *Rarity {

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

	collSize := float64(len(coll))
	var entropy float64
	for attr, dist := range distributions {
		for _, n := range dist {
			p := n / collSize
			entropy += -p * math.Log2(p)
		}
		// null-Value information
		if p := (collSize - counts[attr]) / collSize; p != 0 {
			entropy += -p * math.Log2(p)
		}
	}

	scores := make(map[TokenID]float64)
	for id := range coll {
		// It's important to calculate over all possible attributes, even those
		// that a particular token lacks. Without this, we would favour tokens
		// that simply have more traits.
		for attr, numHaveAttr := range counts {
			var n float64
			if v, ok := attributes[id][attr]; ok {
				n = distributions[attr][v]
			} else {
				n = collSize - numHaveAttr
			}

			scores[id] += -math.Log2(n / collSize)
		}
	}

	// It's not valid to consider all decimal points so we limit the precision
	// relative to the log of the collection size as this is what dictates
	// precision of individual probabilities.
	scale := (func() func(float64) float64 {
		precision := int(math.Floor(math.Log10(float64(len(coll)))))
		pow := math.Pow10(precision)
		return func(f float64) float64 {
			return math.Round(f*pow) / pow
		}
	})()

	for id := range scores {
		scores[id] = scale(scores[id] / entropy)
	}
	return &Rarity{
		Entropy: entropy,
		Scores:  scores,
	}
}
