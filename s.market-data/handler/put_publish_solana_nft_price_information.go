package handler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"
)

type SolanaNFTInfo struct {
	CollectionID          string
	VendorID              string
	Price                 float64
	HumanizedCollectionID string
	Vendor                string
	Emoji                 string
	Price4H               float64
	Price24H              float64
}

var (
	solanaNFTAssets = assets.SolanaNFTAssets
)

var (
	solanaNFTPriceCache4h  *ttlcache.Cache
	solanaNFTPriceCache24h *ttlcache.Cache
	solanaNFTPriceOnce     sync.Once
)

// PublishSolanaNFTPriceInformation ...
func (s *MarketDataService) PublishSolanaNFTPriceInformation(
	ctx context.Context, in *marketdataproto.PublishSolanaNFTPriceInformationRequest,
) (*marketdataproto.PublishSolanaNFTPriceInformationResponse, error) {
	// Init caches.
	solanaNFTPriceOnce.Do(func() {
		now := time.Now().UTC()

		// For offsets see: https://www.notion.so/Solana-NFT-Floor-monitor-1d5b44caf8c54f8d885d71fa0a3f928d
		// Calc offset for 4h.
		truncNow4h := now.Truncate(4 * time.Hour)
		ttl4h := truncNow4h.Add(4 * time.Hour).Sub(now)

		solanaNFTPriceCache4h = ttlcache.NewCache()
		solanaNFTPriceCache4h.SetCacheSizeLimit(len(solanaNFTAssets))
		solanaNFTPriceCache4h.SetTTL(ttl4h)
		solanaNFTPriceCache4h.SkipTTLExtensionOnHit(false)

		// Calc offset for 24h
		truncNow24h := now.Truncate(4 * time.Hour)
		ttl24h := truncNow24h.Add(4 * time.Hour).Sub(now)

		solanaNFTPriceCache24h = ttlcache.NewCache()
		solanaNFTPriceCache24h.SetCacheSizeLimit(len(solanaNFTAssets))
		solanaNFTPriceCache24h.SetTTL(ttl24h)
		solanaNFTPriceCache24h.SkipTTLExtensionOnHit(false)
	})

	slog.Trace(ctx, "Fetching & publishing solana NFT price information for [%d] assets", len(solanaNFTAssets))

	var (
		nfts []*SolanaNFTInfo
		wg   sync.WaitGroup
		mu   sync.Mutex
	)
	for _, nft := range solanaNFTAssets {
		nft := nft
		wg.Add(1)

		go func() {
			defer wg.Done()
			time.Sleep(jitter(0, 60))

			rsp, err := getSolanaNFTFloorPrice(ctx, nft.CollectionID, nft.Vendor)
			if err != nil {
				slog.Error(ctx, "Failed to get floor price for solana NFT: %s from %s: Error: %v", nft.CollectionID, nft.Vendor, err)
				return
			}

			// We have a limit of one so this is what we expect. But lets be defensive.
			var solanaNftInfo *SolanaNFTInfo
			switch len(rsp) {
			case 0:
				return
			default:
			}

			r := rsp[0]
			solanaNftInfo = &SolanaNFTInfo{
				Price:                 float64(r.Price),
				CollectionID:          nft.CollectionID,
				HumanizedCollectionID: nft.HumanizedCollectionID,
				VendorID:              r.Id,
				Vendor:                r.Vendor.String(),
				Emoji:                 nft.Emoji,
			}

			var (
				key               = fmt.Sprintf("%s-%s", nft.CollectionID, nft.Vendor)
				price4h, price24h float64
			)

			// Fetch and / or set 4h cache price.
			v4h, ttl4h, err := solanaNFTPriceCache4h.GetWithTTL(key)
			switch {
			case err == ttlcache.ErrNotFound:
				solanaNFTPriceCache4h.SetWithTTL(key, float64(r.Price), 4*time.Hour)
				price4h = float64(r.Price)
			case err != nil:
				slog.Error(ctx, "Failed to get ttl cache for 4h: %s", key)
			case ttl4h == 0:
				solanaNFTPriceCache4h.SetWithTTL(key, float64(r.Price), 4*time.Hour)
				price4h = float64(r.Price)
			default:
				f, ok := v4h.(float64)
				switch {
				case !ok:
					slog.Error(ctx, "Failed to convert 4h floor price cache value to float: %T", v4h)
					price4h = float64(r.Price)
				default:
					price4h = f
				}
			}

			// Fetch and / or set 24h cache price.
			v24h, ttl24h, err := solanaNFTPriceCache24h.GetWithTTL(key)
			switch {
			case err == ttlcache.ErrNotFound:
				solanaNFTPriceCache24h.SetWithTTL(key, float64(r.Price), 4*time.Hour)
				price24h = float64(r.Price)
			case err != nil:
				slog.Error(ctx, "Failed to get ttl cache for 24h: %s", key)
			case ttl24h == 0:
				solanaNFTPriceCache24h.SetWithTTL(key, float64(r.Price), 4*time.Hour)
				price24h = float64(r.Price)
			default:
				f, ok := v24h.(float64)
				switch {
				case !ok:
					slog.Error(ctx, "Failed to convert 24h floor price cache value to float: %T", v24h)
					price24h = float64(r.Price)
				default:
					price24h = f
				}
			}

			// Set prices from cache.
			solanaNftInfo.Price4H = price4h
			solanaNftInfo.Price24H = price24h

			mu.Lock()
			defer mu.Unlock()

			nfts = append(nfts, solanaNftInfo)
		}()
	}

	// Wait for all goroutines to complete.
	wg.Wait()

	// Sort based on collection ID.
	sort.Slice(nfts, func(i, j int) bool {
		ni, nj := nfts[i], nfts[j]

		var ki, kj string
		switch {
		case ni.CollectionID != "" && nj.CollectionID != "":
			ki, kj = ni.HumanizedCollectionID, nj.HumanizedCollectionID
		default:
			ki, kj = ni.CollectionID, nj.CollectionID
		}

		return ki < kj
	})

	// Calc. collection indent.
	var collectionIDIndent int
	for _, nft := range nfts {
		if len(nft.HumanizedCollectionID) > collectionIDIndent {
			collectionIDIndent = len(nft.HumanizedCollectionID)
		}
	}

	// Calc. vendor indent.
	var vendorIndent int
	for _, nft := range nfts {
		if len(nft.Vendor) > vendorIndent {
			vendorIndent = len(nft.Vendor)
		}
	}

	now := time.Now().UTC().Truncate(time.Hour)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n:robot:    `Solana NFTs Floor Price Data [%v]`     :dove:\n", now))

	for _, nft := range nfts {
		priceDiff4h := nft.Price - nft.Price4H
		var emoji4h string
		switch {
		case priceDiff4h > 0:
			emoji4h = ":small_red_triangle:"
		default:
			emoji4h = ":small_red_triangle_down:"
		}

		priceDiff24h := nft.Price - nft.Price24H
		var emoji24h string
		switch {
		case priceDiff24h > 0:
			emoji24h = ":small_red_triangle:"
		default:
			emoji24h = ":small_red_triangle_down:"
		}

		sb.WriteString(
			fmt.Sprintf(
				"\n%s` [%s]:%s %s%s CURRENT: %.2f SOL 4h: %.2f SOL` %s `24h: %.2f SOL` %s `# LISTED: %d`",
				nft.Emoji,
				nft.HumanizedCollectionID,
				addPadding(collectionIDIndent-len(nft.HumanizedCollectionID)+1),
				nft.Vendor,
				addPadding(vendorIndent-len(nft.Vendor)+1),
				nft.Price,
				nft.Price4H,
				emoji4h,
				nft.Price24H,
				emoji24h,
				len(nfts),
			),
		)
	}

	idempotencyKey := fmt.Sprintf("solananftfloorprice-%s", now)
	if err := publishToDiscord(ctx, sb.String(), discordproto.DiscordSatoshiNFTBotChannel, idempotencyKey); err != nil {
		return nil, gerrors.Augment(err, "failed_to_publish_solana_nft_price_information", map[string]string{
			"idempotency_key": idempotencyKey,
		})
	}

	return &marketdataproto.PublishSolanaNFTPriceInformationResponse{}, nil
}
