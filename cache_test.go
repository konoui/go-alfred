package alfred

import (
	"testing"
	"time"
)

func TestCache_Store(t *testing.T) {
	tests := []struct {
		name string

		wantErr bool
	}{
		{name: "store", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wf := testWorkflow().Append(NewItem().Title("title").Subtitle("subtitle"))
			c := wf.Cache("test")
			if err := c.Store(); (err != nil) != tt.wantErr {
				t.Errorf("Cache.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_Load(t *testing.T) {
	tests := []struct {
		name string

		wantErr bool
	}{
		{name: "load", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "test"
			wf := testWorkflow().Append(NewItem().Title("title").Subtitle("subtitle"))
			if err := wf.Cache(key).Store(); err != nil {
				t.Fatalf("Cache.StoreItems() error = %v", err)
			}

			want := wf.Bytes()

			gwf := testWorkflow()
			if err := gwf.Cache(key).MaxAge(60 * time.Second).Load(); err != nil {
				t.Errorf("Cache.Load() error = %v", err)
			}

			got := gwf.Bytes()
			if diff := DiffOutput(want, got); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}
		})
	}
}
