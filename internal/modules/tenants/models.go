package tenants

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	Slug       string     `db:"slug" json:"slug"`
	SchemaName string     `db:"schema_name" json:"schema_name"`
	Domain     *string    `db:"domain" json:"domain,omitempty"`
	Status     string     `db:"status" json:"status"`
	PlanID     *uuid.UUID `db:"plan_id" json:"plan_id,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateTenantRequest struct {
	Name   string  `json:"name" validate:"required"`
	Slug   string  `json:"slug" validate:"required"`
	Domain *string `json:"domain,omitempty"`
	PlanID *string `json:"plan_id,omitempty"`
}

type UpdateTenantRequest struct {
	Name   *string `json:"name,omitempty"`
	Domain *string `json:"domain,omitempty"`
	Status *string `json:"status,omitempty"`
	PlanID *string `json:"plan_id,omitempty"`
}

type TenantResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	SchemaName string    `json:"schema_name"`
	Domain     *string   `json:"domain,omitempty"`
	Status     string    `json:"status"`
	PlanID     *string   `json:"plan_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (t *Tenant) ToResponse() *TenantResponse {
	resp := &TenantResponse{
		ID:         t.ID.String(),
		Name:       t.Name,
		Slug:       t.Slug,
		SchemaName: t.SchemaName,
		Domain:     t.Domain,
		Status:     t.Status,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}

	if t.PlanID != nil {
		planID := t.PlanID.String()
		resp.PlanID = &planID
	}

	return resp
}
