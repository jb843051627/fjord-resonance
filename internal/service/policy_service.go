package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type PolicyService struct {
	mu    sync.RWMutex
	rules []model.AlertRule
}

func NewPolicyService() *PolicyService { return &PolicyService{rules: model.DefaultAlertRules()} }

func (s *PolicyService) Rules() []model.AlertRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.AlertRule(nil), s.rules...)
}

func (s *PolicyService) Replace(ctx context.Context, rules []model.AlertRule) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(rules) == 0 {
		return fmt.Errorf("policy list is empty")
	}
	s.mu.Lock()
	s.rules = append([]model.AlertRule(nil), rules...)
	s.mu.Unlock()
	return nil
}

func (s *PolicyService) Match(alert model.Alert) (model.AlertRule, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return model.MatchRule(alert, s.rules)
}
