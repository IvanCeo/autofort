package entity_test

import (
	"errors"
	"testing"

	"autofort/internal/entity"

	v10 "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type wantFieldTag struct {
	field string
	tag   string
}

func TestCustomerValidate(t *testing.T) {
	tests := []struct {
		name string
		c    entity.Customer
		// wantNil=true => ошибок быть не должно
		wantNil bool
		// минимальные ожидания по (field, tag)
		want []wantFieldTag
	}{
		{
			name: "ok",
			c: entity.Customer{
				ID:          uuid.New(),
				FirstName:   "Ivan",
				LastName:    "Mikhaylov",
				PhoneNumber: "+79272110520",
			},
			wantNil: true,
		},
		{
			name: "8-type phone number",
			c: entity.Customer{
				ID:          uuid.New(),
				FirstName:   "Ivan",
				LastName:    "Mi",
				PhoneNumber: "89272110520",
			},
			wantNil: true,
		},
		{
			name: "missing all required fields",
			c:    entity.Customer{},
			want: []wantFieldTag{
				{"FirstName", "required"},
				{"LastName", "required"},
				{"PhoneNumber", "required"},
			},
		},
		{
			name: "invalid phone e164",
			c: entity.Customer{
				ID:          uuid.New(),
				FirstName:   "Ivan",
				LastName:    "Mikhaylov",
				PhoneNumber: "123",
			},
			want: []wantFieldTag{
				{"PhoneNumber", "e164"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()

			if tt.wantNil {
				if err != nil {
					t.Fatalf("expected nil error, got: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected validation error, got nil")
			}

			var ve v10.ValidationErrors
			if !errors.As(err, &ve) {
				t.Fatalf("expected validator.ValidationErrors, got %T: %v", err, err)
			}

			for _, w := range tt.want {
				if !hasFieldTag(ve, w.field, w.tag) {
					t.Fatalf("expected field=%s tag=%s, got: %s", w.field, w.tag, summarize(ve))
				}
			}
		})
	}
}

func hasFieldTag(ve v10.ValidationErrors, field, tag string) bool {
	for _, fe := range ve {
		if fe.Field() == field && fe.Tag() == tag {
			return true
		}
	}
	return false
}

func summarize(ve v10.ValidationErrors) string {
	out := ""
	for i, fe := range ve {
		if i > 0 {
			out += ", "
		}
		out += fe.Field() + ":" + fe.Tag()
	}
	return out
}
