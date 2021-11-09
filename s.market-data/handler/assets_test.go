package handler

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"swallowtail/s.market-data/assets"
)

func TestAssets_Sort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		unsortedAssets       AssetInfoList
		expectedSortedAssets AssetInfoList
	}{
		{
			name: "basic_same_group",
			unsortedAssets: AssetInfoList{
				{
					Symbol: "ETH",
					Group:  assets.AssetGroupBitcoin.String(),
				},
				{
					Symbol: "BTC",
					Group:  assets.AssetGroupBitcoin.String(),
				},
			},
			expectedSortedAssets: AssetInfoList{
				{
					Symbol: "BTC",
					Group:  assets.AssetGroupBitcoin.String(),
				},
				{
					Symbol: "ETH",
					Group:  assets.AssetGroupBitcoin.String(),
				},
			},
		},
		{
			name: "basic_diff_groups",
			unsortedAssets: AssetInfoList{
				{
					Symbol: "CRV",
					Group:  assets.AssetGroupDeFi.String(),
				},
				{
					Symbol: "BTC",
					Group:  assets.AssetGroupBitcoin.String(),
				},
				{
					Symbol: "AAVE",
					Group:  assets.AssetGroupDeFi.String(),
				},
			},
			expectedSortedAssets: AssetInfoList{
				{
					Symbol: "BTC",
					Group:  assets.AssetGroupBitcoin.String(),
				},
				{
					Symbol: "AAVE",
					Group:  assets.AssetGroupDeFi.String(),
				},
				{
					Symbol: "CRV",
					Group:  assets.AssetGroupDeFi.String(),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sort.Sort(tt.unsortedAssets)

			for _, a := range tt.unsortedAssets {
				fmt.Printf("\n%+v", a)
			}

			assert.Equal(t, tt.expectedSortedAssets, tt.unsortedAssets)
		})
	}
}
