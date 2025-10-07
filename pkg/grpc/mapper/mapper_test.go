package mapper_test

import (
	"gophkeeper/pkg/grpc/auth"
	"gophkeeper/pkg/grpc/mapper"
	"reflect"
	"testing"
)

func TestKdfParameters_RoundTrip(t *testing.T) {
	orig := &auth.KdfParams{
		Alg:         auth.KdfAlg_ARGON2ID,
		TimeCost:    3,
		MemoryCost:  65536,
		Parallelism: 4,
		Salt:        []byte("salt123"),
	}

	// gRPC -> domain -> gRPC
	domain := mapper.KdfParametersToDomain(orig)
	back := mapper.KdfParametersToGRPC(domain)

	if !reflect.DeepEqual(orig, back) {
		t.Errorf("roundtrip mismatch:\norig=%+v\nback=%+v", orig, back)
	}
}

func TestAuthAlgo_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		in   auth.AuthKeyAlg
	}{
		{"HMAC_SHA256", auth.AuthKeyAlg_HMAC_SHA256},
		{"HMAC_SHA512", auth.AuthKeyAlg_HMAC_SHA512},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mapper.AuthAlgoToDomain(tt.in)
			back := mapper.AuthAlgoToGRPC(d)
			if back != tt.in {
				t.Errorf("roundtrip mismatch: got %v, want %v", back, tt.in)
			}
		})
	}
}

func TestKdfAlgoToDomain_PanicOnUnknown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on unknown KDF alg, got none")
		}
	}()
	_ = mapper.KdfParametersToDomain(&auth.KdfParams{
		Alg: auth.KdfAlg(999),
	})
}
