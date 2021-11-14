package alfred

import (
	"os"
	"testing"
)

func Benchmark_initAssets(b *testing.B) {
	w := NewWorkflow(
		WithLogWriter(os.Stderr),
	)

	if err := os.RemoveAll(w.getAssetsDir()); err != nil {
		b.Fatal(err)
	}

	if err := new(envs).Initialize(w); err != nil {
		b.Fatal(err)
	}

	if err := new(assets).Initialize(w); err != nil {
		b.Fatal(err)
	}
}
