package policy

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Policy, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]Policy, error)
	Create(ctx context.Context, org_id string, input CreatePolicyRequest) (Policy, error)
	Update(ctx context.Context, org_id string, id string, input UpdatePolicyRequest) (Policy, error)
	Patch(ctx context.Context, org_id string, id string, input PatchPolicyRequest) (Policy, error)
	Delete(ctx context.Context, org_id string, id string) error
	// Patch(ctx context.Context, org_id string, id string, req UserPatchRequest) (User, error)
}

type Policy struct {
	mongo_entity.Policy
}

type CreatePolicyRequest struct {
	Identifier  string `json:"identifier" bson:"identifier"`
	DisplayName string `json:"display_name" bson:"display_name"`
	Version     string `json:"version" bson:"version"`
	Policy      string `json:"policy" bson:"policy"`
}

type UpdatePolicyRequest struct {
	DisplayName   *string              `json:"display_name,omitempty" bson:"display_name"`
	ActiveVersion *string              `json:"active_version,omitempty" bson:"active_version"`
	PolicyContent *UpdatePolicyContent `json:"policy_content,omitempty" bson:"policy_content"`
}

type PatchPolicyRequest struct {
	AddedPolicies   []mongo_entity.PolicyContent `json:"added_policies,omitempty" bson:"added_policies"`
	RemovedPolicies []string                     `json:"removed_policies,omitempty" bson:"removed_policies"`
}

type UpdatePolicy struct {
	DisplayName   *string              `json:"display_name,omitempty" bson:"display_name"`
	ActiveVersion *string              `json:"active_version,omitempty" bson:"active_version"`
	PolicyContent *UpdatePolicyContent `json:"policy_content,omitempty" bson:"policy_content"`
}

type UpdatePolicyContent struct {
	Version *string `json:"version" bson:"version"`
	Policy  *string `json:"policy" bson:"policy"`
}

type PatchPolicy struct {
	AddedPolicies   []mongo_entity.PolicyContent `json:"added_policies,omitempty" bson:"added_policies"`
	RemovedPolicies []string                     `json:"removed_policies,omitempty" bson:"removed_policies"`
}

func (m CreatePolicyRequest) Validate() error {

	return validation.ValidateStruct(&m,
		validation.Field(&m.Identifier, validation.Required),
		validation.Field(&m.Version, validation.Required),
		validation.Field(&m.Policy, validation.Required),
	)
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {

	return service{repo: repo, logger: logger}
}

// Get policy by id.
func (s service) Get(ctx context.Context, org_id string, id string) (Policy, error) {

	policy, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Error while getting the user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return Policy{}, &util.NotFoundError{Path: "User"}
	}
	return Policy{*policy}, err
}

// Create new policy.
func (s service) Create(ctx context.Context, org_id string, req CreatePolicyRequest) (Policy, error) {

	// Validate policy request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating policy create request.")
		return Policy{}, &util.InvalidInputError{Path: "Invalid input for policy."}
	}

	// Check policy already exists.
	exists, _ := s.repo.CheckPolicyExistsByIdentifier(ctx, org_id, req.Identifier)
	if exists {
		s.logger.Debug("User already exists.")
		return Policy{}, &util.AlreadyExistsError{Path: "User : " + req.Identifier}

	}

	// Generate policy id.
	policyId := primitive.NewObjectID()
	policContentId := primitive.NewObjectID()
	policyContent := mongo_entity.PolicyContent{
		ID:      policContentId,
		Version: req.Version,
		Policy:  req.Policy,
	}
	err := s.repo.Create(ctx, org_id, mongo_entity.Policy{
		ID:             policyId,
		DisplayName:    req.DisplayName,
		Identifier:     req.Identifier,
		ActiveVersion:  req.Version,
		PolicyContents: []mongo_entity.PolicyContent{policyContent},
	})

	if err != nil {
		s.logger.Error("Error while creating user.",
			zap.String("organization_id", org_id))
		return Policy{}, err
	}
	return s.Get(ctx, org_id, policyId.Hex())
}

// // Update policy.
func (s service) Update(ctx context.Context, org_id string, id string, req UpdatePolicyRequest) (Policy, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Policy not exists.", zap.String("policy_id", id))
		return Policy{}, &util.NotFoundError{Path: "Policy " + id + " not exists."}
	}

	if req.PolicyContent != nil {
		if req.PolicyContent.Version == nil || *req.PolicyContent.Version == "" {
			return Policy{}, &util.InvalidInputError{Path: "Invalid input for policy."}
		}
		exists, _ := s.repo.CheckPolicyContentExistsByVersion(ctx, org_id, *req.PolicyContent.Version)
		if !exists {
			return Policy{}, &util.InvalidInputError{Path: "Invalid policy version " + *req.PolicyContent.Version}
		}
		if req.PolicyContent.Policy == nil || *req.PolicyContent.Policy == "" {
			return Policy{}, &util.InvalidInputError{Path: "Invalid input for policy."}
		}
	}

	if err := s.repo.Update(ctx, org_id, id, UpdatePolicy{
		DisplayName:   req.DisplayName,
		ActiveVersion: req.ActiveVersion,
		PolicyContent: req.PolicyContent,
	}); err != nil {
		s.logger.Error("Error while updating user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return Policy{}, err
	}
	updatedPolicy, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("User not exists.", zap.String("user_id", id))
		return Policy{}, &util.NotFoundError{Path: "User " + id + " not exists."}
	}
	return Policy{*updatedPolicy}, nil
}

func (s service) Patch(ctx context.Context, org_id string, id string, req PatchPolicyRequest) (Policy, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("User not exists.", zap.String("user_id", id))
		return Policy{}, &util.NotFoundError{Path: "User " + id + " not exists."}
	}

	// roles
	for _, policy := range req.AddedPolicies {
		exists, _ := s.repo.CheckPolicyContentExistsByVersion(ctx, org_id, policy.Version)
		if exists {
			return Policy{}, &util.InvalidInputError{Path: "Invalid policy version " + policy.Version}
		}
	}
	for _, version := range req.RemovedPolicies {
		exists, _ := s.repo.CheckPolicyContentExistsByVersion(ctx, org_id, version)
		if !exists {
			return Policy{}, &util.InvalidInputError{Path: "Invalid policy version " + version}
		}
	}

	added_policies := []mongo_entity.PolicyContent{}
	for _, policyContent := range req.AddedPolicies {
		policyContentId := primitive.NewObjectID()
		added_policies = append(added_policies, mongo_entity.PolicyContent{
			ID:      policyContentId,
			Version: policyContent.Version,
			Policy:  policyContent.Policy,
		})
	}
	if err := s.repo.Patch(ctx, org_id, id, PatchPolicy{
		AddedPolicies:   added_policies,
		RemovedPolicies: req.RemovedPolicies,
	}); err != nil {
		s.logger.Error("Error while updating user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return Policy{}, err
	}
	updatedUser, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("User not exists.", zap.String("user_id", id))
		return Policy{}, &util.NotFoundError{Path: "User " + id + " not exists."}
	}
	return Policy{*updatedUser}, nil
}

// Delete user.
func (s service) Delete(ctx context.Context, org_id string, id string) error {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("User not exists.", zap.String("user_id", id))
		return &util.NotFoundError{Path: "User " + id + " not exists."}

	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		s.logger.Error("Error while deleting user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return err
	}
	return nil
}

// Pagination filter.
type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

// // Get all user.
func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]Policy, error) {

	result := []Policy{}
	items, err := s.repo.Query(ctx, org_id)
	if err != nil {
		s.logger.Error("Error while retrieving all resources.",
			zap.String("organization_id", org_id))
		return []Policy{}, err
	}

	for _, item := range *items {
		result = append(result, Policy{item})
	}
	return result, err
}
