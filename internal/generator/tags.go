package generator

import "fmt"
import "ar/internal/generator/rand"

// generate random account and instance key:value "tags"
func randomTags() []string {
	t := []string{"env:local"}
	t = append(t, fmt.Sprintf("account:%d", rand.SeededRand.Int()%10000))
	t = append(t, "instance:"+rand.String(10, rand.CharsetLower))
	return t
}
