package auth

import "testing"

func TestHashPassword(t *testing.T){
	hash, err := HashPassword("testing")
    if err != nil {
        t.Fatalf("HashPassword returned unexpected error: %v", err)
    }

	isEqual, err := CheckPasswordHash("testing", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned unexpected error: %v", err)
	}

    if !isEqual {
        t.Error("expected hash to match original password, but it did not")
    }
}
