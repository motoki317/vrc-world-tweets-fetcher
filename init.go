package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/sivchari/gotwtr"

	"github.com/motoki317/vrc-world-tweets-fetcher/api/twitter"
)

const targetRuleValue = "(#MadeWithVRChat OR #VRChat_world OR #VRChat_world紹介 OR #VRChatワールド紹介) -is:retweet -is:reply -is:quote"

var targetRule = &gotwtr.AddRule{Value: targetRuleValue, Tag: ""}

func cmdInit() error {
	log.Printf("Syncing stream rules. Target rule value: %s\n", targetRuleValue)

	log.Println("Fetching current rules...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rules, err := twitter.FetchStreamRules(ctx)
	if err != nil {
		return fmt.Errorf("encountered an error while fetching existing stream rules: %w", err)
	}
	if len(rules.Errors) > 0 {
		return fmt.Errorf("twitter API returned an error while fetching existing stream rules: %s", rules.Errors[0].Message)
	}

	log.Printf("There are currently %d rule(s) defined.\n", len(rules.Rules))
	for i, r := range rules.Rules {
		log.Printf("%d. value: %s, tag: %s, rule_id: %s", i+1, r.Value, r.Tag, r.ID)
	}
	matchingRuleIDs := lo.Map[*gotwtr.FilteredRule, string](
		lo.Filter[*gotwtr.FilteredRule](rules.Rules, func(r *gotwtr.FilteredRule, _ int) bool { return r.Value == targetRuleValue }),
		func(r *gotwtr.FilteredRule, _ int) string { return r.ID },
	)
	nonMatchingRuleIDs := lo.Map[*gotwtr.FilteredRule, string](
		lo.Filter[*gotwtr.FilteredRule](rules.Rules, func(r *gotwtr.FilteredRule, _ int) bool { return r.Value != targetRuleValue }),
		func(r *gotwtr.FilteredRule, _ int) string { return r.ID },
	)
	log.Printf("Of which, %d rule(s) match the target rule value (%s).\n", len(matchingRuleIDs), strings.Join(matchingRuleIDs, ", "))

	switch len(matchingRuleIDs) {
	case 0:
		log.Printf("Found no matching rules. Will add 1 target rule, and delete %d existing rule(s).\n", len(nonMatchingRuleIDs))
		if res, err := twitter.DeleteStreamRules(ctx, nonMatchingRuleIDs); err != nil {
			if len(res.Errors) > 0 {
				log.Printf("delete error: %s\n", res.Errors[0].Message)
			}
			return fmt.Errorf("encountered an error while adding and deleting rules: %w", err)
		}
		if res, err := twitter.AddStreamRules(ctx, []*gotwtr.AddRule{targetRule}); err != nil {
			if len(res.Errors) > 0 {
				log.Printf("add error: %s\n", res.Errors[0].Message)
			}
			return fmt.Errorf("encountered an error while adding and deleting rules: %w", err)
		}
	case 1:
		if len(rules.Rules) == 1 {
			log.Println("Found exactly 1 rule and the rule was the desired rule! No operations required.")
			return nil
		}
		log.Printf("Found exactly 1 matching rule. Will delete all %d non-matching rule(s).\n", len(nonMatchingRuleIDs))
		if res, err := twitter.DeleteStreamRules(ctx, nonMatchingRuleIDs); err != nil {
			if len(res.Errors) > 0 {
				log.Printf("delete error: %s\n", res.Errors[0].Message)
			}
			return fmt.Errorf("encountered an error while deleting non-matching rules: %w", err)
		}
	default:
		log.Printf("Found %d matching rules. Will delete unnecessary %d rule(s) and %d non-matching rule(s).\n", len(matchingRuleIDs), len(matchingRuleIDs)-1, len(nonMatchingRuleIDs))
		deleteIDs := append(append([]string{}, nonMatchingRuleIDs...), matchingRuleIDs[1:]...)
		if res, err := twitter.DeleteStreamRules(ctx, deleteIDs); err != nil {
			if len(res.Errors) > 0 {
				log.Printf("delete error: %s\n", res.Errors[0].Message)
			}
			return fmt.Errorf("encountered an error while deleting non-matching rules: %w", err)
		}
	}

	log.Println("Successfully synced rule values!")
	return nil
}
