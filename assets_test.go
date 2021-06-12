package alfred

import (
	"os"
	"testing"
)

func Benchmark_initAssets(b *testing.B) {
	w := NewWorkflow()
	w.SetLog(os.Stderr)
	if err := os.RemoveAll(w.getAssetsDir()); err != nil {
		b.Fatal(err)
	}

	if err := w.init(); err != nil {
		b.Fatal(err)
	}

	if err := w.initAssets(); err != nil {
		b.Fatal(err)
	}
}
