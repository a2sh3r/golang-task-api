package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantError bool
	}{
		{
			name:      "valid debug",
			level:     "debug",
			wantError: false,
		},
		{
			name:      "valid info",
			level:     "info",
			wantError: false,
		},
		{
			name:      "invalid level",
			level:     "abracadabra",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.level)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, Log)
				// Проверим, что Log не nil и не zap.NewNop()
				assert.IsType(t, &zap.Logger{}, Log)
			}
		})
	}
}
