package attacker

import (
	"context"
	"testing"

	vegeta "github.com/NoneTirex/vegeta/v12/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/NoneTirex/ali/storage"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNewAttacker(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		opts    Options
		wantErr bool
	}{
		{
			name:    "no target given",
			target:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAttacker(&storage.FakeStorage{}, tt.target, &tt.opts)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAttack(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name    string
		target  string
		opts    Options
		wantErr bool
	}{
		{
			name:   "no result given back",
			target: "http://host.xz",
			opts: Options{
				Attacker: &fakeBackedAttacker{},
			},
		},
		{
			name:   "two result given back",
			target: "http://host.xz",
			opts: Options{
				Attacker: &fakeBackedAttacker{
					results: []*vegeta.Result{
						{
							Code: 200,
						},
						{
							Code: 200,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAttacker(&storage.FakeStorage{}, tt.target, &tt.opts)
			require.NoError(t, err)
			metricsCh := make(chan *Metrics, 100)
			a.Attack(ctx, metricsCh)
		})
	}
}
