package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVRChatWorldIDs(t *testing.T) {
	tests := []struct {
		name string
		urls []string
		want []VRChatWorldID
	}{
		{
			name: "single URL 1",
			urls: []string{"https://vrchat.com/home/world/wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
			want: []VRChatWorldID{"wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
		},
		{
			name: "single URL 1 alt",
			urls: []string{"https://vrchat.com/home/world/wrld_fa576623-79cd-412a-ad99-9af4870313df"},
			want: []VRChatWorldID{"wrld_fa576623-79cd-412a-ad99-9af4870313df"},
		},
		{
			name: "single URL 2",
			urls: []string{"https://vrchat.com/home/launch?worldId=wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
			want: []VRChatWorldID{"wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
		},
		{
			name: "with query 1",
			urls: []string{"https://vrchat.com/home/world/wrld_c4e61164-c854-42d3-a70c-17d0d68a141b?someQueryParam=1"},
			want: []VRChatWorldID{"wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
		},
		{
			name: "with query 2",
			urls: []string{"https://vrchat.com/home/launch?worldId=wrld_c4e61164-c854-42d3-a70c-17d0d68a141b&someQueryParam=1"},
			want: []VRChatWorldID{"wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
		},
		{
			name: "with query 3",
			urls: []string{"https://vrchat.com/home/launch?someQueryParam=1&worldId=wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
			want: []VRChatWorldID{"wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
		},
		{
			name: "mixed URLs",
			urls: []string{
				"https://example.com",
				"https://vrchat.com/home/world/wrld_c4e61164-c854-42d3-a70c-17d0d68a141b",
			},
			want: []VRChatWorldID{"wrld_c4e61164-c854-42d3-a70c-17d0d68a141b"},
		},
		{
			name: "invalid single URL 1",
			urls: []string{"https://vrchat.com/home/world/world_c4e61164-c854-42d3-a70c-17d0d68a141b"}, // world_... instead of wrld_...
			want: []VRChatWorldID{},
		},
		{
			name: "invalid single URL 2",
			urls: []string{"https://vrchat.com/home/launch?worldId=world_c4e61164-c854-42d3-a70c-17d0d68a141b"},
			want: []VRChatWorldID{},
		},
		{
			name: "invalid with query 1",
			urls: []string{"https://vrchat.com/home/world/world_c4e61164-c854-42d3-a70c-17d0d68a141b?someQueryParam=1"},
			want: []VRChatWorldID{},
		},
		{
			name: "invalid with query 2",
			urls: []string{"https://vrchat.com/home/launch?worldId=world_c4e61164-c854-42d3-a70c-17d0d68a141b&someQueryParam=1"},
			want: []VRChatWorldID{},
		},
		{
			name: "invalid with query 3",
			urls: []string{"https://vrchat.com/home/launch?someQueryParam=1&worldId=world_c4e61164-c854-42d3-a70c-17d0d68a141b"},
			want: []VRChatWorldID{},
		},
		{
			name: "invalid mixed URLs",
			urls: []string{
				"https://example.com",
				"https://vrchat.com/home/world/world_c4e61164-c854-42d3-a70c-17d0d68a141b",
			},
			want: []VRChatWorldID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractVRChatWorldIDs(tt.urls); !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("ExtractVRChatWorldIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
