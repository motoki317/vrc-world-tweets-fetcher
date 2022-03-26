package twitter

import (
	"context"

	"github.com/sivchari/gotwtr"
)

func FetchStreamRules(ctx context.Context) (*gotwtr.RetrieveStreamRulesResponse, error) {
	return c.RetrieveStreamRules(ctx)
}

func AddStreamRules(ctx context.Context, addRules []*gotwtr.AddRule) (*gotwtr.AddOrDeleteRulesResponse, error) {
	return c.AddOrDeleteRules(ctx, &gotwtr.AddOrDeleteJSONBody{
		Add: addRules,
	})
}

func DeleteStreamRules(ctx context.Context, deleteRuleIDs []string) (*gotwtr.AddOrDeleteRulesResponse, error) {
	return c.AddOrDeleteRules(ctx, &gotwtr.AddOrDeleteJSONBody{
		Delete: &gotwtr.DeleteRule{IDs: deleteRuleIDs},
	})
}
