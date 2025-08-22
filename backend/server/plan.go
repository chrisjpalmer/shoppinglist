package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/gen"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
)

func (s *Server) GetPlan(ctx context.Context, rq *connect.Request[gen.GetPlanRequest]) (*connect.Response[gen.GetPlanResponse], error) {
	p, err := s.db.GetPlan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		gp, err := s.createEmptyPlan(ctx)
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&gen.GetPlanResponse{Plan: &gp}), nil
	}

	if err != nil {
		return nil, err
	}

	var gp gen.Plan
	err = unmarshalJSON(p.PlanData, &gp)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GetPlanResponse{Plan: &gp}), nil
}

func (s *Server) UpdatePlan(ctx context.Context, rq *connect.Request[gen.UpdatePlanRequest]) (*connect.Response[gen.UpdatePlanResponse], error) {
	p, err := s.db.GetPlan(ctx)

	pstr, err := marshalJSON(rq.Msg.Plan)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		if err := s.db.CreatePlan(ctx, pstr); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	err = s.db.UpdatePlan(ctx, generated.UpdatePlanParams{
		ID:       p.ID,
		PlanData: pstr,
	})

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.UpdatePlanResponse{}), nil
}

func (s *Server) createEmptyPlan(ctx context.Context) (gen.Plan, error) {
	p := emptyPlan()

	ps, err := marshalJSON(p)
	if err != nil {
		return gen.Plan{}, err
	}

	err = s.db.CreatePlan(ctx, ps)
	if err != nil {
		return gen.Plan{}, err
	}

	return p, nil
}

func emptyPlan() gen.Plan {
	var days []*gen.Day
	for range 7 {
		days = append(days, &gen.Day{
			CategoryMeals: []*gen.CategoryMeal{
				// 0 = lunch, 1 = dinner, 2 = snack
				{CategoryId: 0}, {CategoryId: 1}, {CategoryId: 2},
			},
		})
	}
	return gen.Plan{
		Days: days,
	}
}

func marshalJSON(obj any) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unmarshalJSON(s string, obj any) error {
	return json.Unmarshal([]byte(s), obj)
}
