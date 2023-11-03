package elastic

import (
	"bytes"
	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/collection/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
	"testing"
)

func TestCollect(t *testing.T) {
	c := collection.New(&bytes.Buffer{})
	Collect(c)
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/elasticsearch/elasticsearch.yml")
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/elasticsearch/elasticsearch.yaml")

	util.AssertAllObfuscatorsTested(t, obfuscators)
}
