package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
)

func hashAddresses(addresses []string) string {
	sorted := append([]string(nil), addresses...)
	sort.Strings(sorted)
	j, _ := json.Marshal(sorted)
	h := sha256.Sum256(j)
	return fmt.Sprintf("%x", h[:])
}
