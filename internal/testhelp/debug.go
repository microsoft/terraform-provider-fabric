// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = (*debugPlan)(nil)

type debugPlan struct{}

func DebugPlan() plancheck.PlanCheck {
	return debugPlan{}
}

func (e debugPlan) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, _ *plancheck.CheckPlanResponse) {
	rd, err := json.Marshal(req.Plan)
	if err != nil {
		_ = fmt.Sprintf("error marshalling machine-readable plan output: %s", err)
	}

	_ = fmt.Sprintf("req.Plan - %s\n", string(rd))
}

// usage in TestStep:
// import "github.com/microsoft/terraform-provider-fabric/internal/testhelp"
//
// ConfigPlanChecks: resource.ConfigPlanChecks{
// 	PostApplyPreRefresh: []plancheck.PlanCheck{
// 		helpers.DebugPlan(),
// 	},
