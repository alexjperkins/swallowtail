package commands

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"

	"swallowtail/libraries/gerrors"
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

const (
	nftCommandID = "nft"
	nftUsage     = `!nft <subcommand>`
)

var (
	collections = map[string]string{
		"babyapes":  solananftsproto.SolanartCollectionIDBabyApes,
		"geckos":    solananftsproto.SolanartCollectionIDGalacticGeckoSpaceGarage,
		"daa":       solananftsproto.SolanartCollectionIDDegenerateApeAcademy,
		"thugbirdz": solananftsproto.SolanartCollectionIDThugBirdz,
	}
)

func init() {
	register(nftCommandID, &Command{
		ID:                  nftCommandID,
		IsPrivate:           false,
		MinimumNumberOfArgs: 1,
		Usage:               nftUsage,
		Description:         "Command for nft data",
		Handler:             nftHandler,
		SubCommands: map[string]*Command{
			"scatter": {
				ID:                  "nft-scatter",
				MinimumNumberOfArgs: 2,
				Usage:               "!nft scatter <collection> <vendor>",
				Description:         "Prints a scattergraph of all nfts in a collection",
				Handler:             scattergraphNFTHandler,
				FailureMsg:          "Please check the vendor & the collection id are correct. No spaces!",
			},
		},
	})
}

func nftHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return gerrors.Unimplemented("parent_command_unimplemented.nft", nil)
}

func scattergraphNFTHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	cID, v := tokens[0], tokens[1]

	// Parse fields.
	vendor, err := convertVendor(v)
	if err != nil {
		return gerrors.Augment(err, "Failed to create scattergraph; bad vendor", nil)
	}

	collectionID, err := convertCollectionID(cID)
	if err != nil {
		return gerrors.Augment(err, "Failed to create scattergraph; bad collection", nil)
	}

	// Gather collection.
	rsp, err := (&solananftsproto.ReadSolanaPriceStatisticsByCollectionIDRequest{
		CollectionId:  collectionID,
		Vendor:        vendor,
		SearchContext: solananftsproto.SearchContextMarketData,
		Limit:         50,
	}).Send(ctx).Response()

	items := rsp.VendorStats

	// Create scattergraph chart for the collection.
	xs := make([]float64, 0, len(items))
	xt := make([]chart.Tick, 0, len(items))

	ys := make([]float64, 0, len(items))
	yt := make([]chart.Tick, 0, len(items))

	annotations := make([]chart.Value2, 0, len(items))

	min, max := math.MaxFloat64, -math.MaxFloat64
	for i, item := range items {
		if float64(item.Price) < min {
			min = float64(item.Price)
		}

		if float64(item.Price) > max {
			max = float64(item.Price)
		}

		xs = append(xs, float64(item.Price))
		xt = append(xt, chart.Tick{
			Value: float64(item.Price),
			Label: item.Name,
		})

		ys = append(ys, float64(i))
		yt = append(yt, chart.Tick{
			Value: float64(i),
			Label: item.Name,
		})

		annotations = append(annotations, chart.Value2{
			XValue: float64(item.Price),
			YValue: float64(i),
			Label:  fmt.Sprintf("%.2f #%d", item.Price, i),
		})
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Price",
			Range: &chart.ContinuousRange{
				Min: min,
				Max: max,
			},
			Ticks:        xt,
			TickPosition: 2,
		},
		YAxis: chart.YAxis{
			Name: "Item Number",
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: float64(len(items)),
			},
			Ticks: yt,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xs,
				YValues: ys,
				Style: chart.Style{
					StrokeWidth:      chart.Disabled,
					Show:             true,
					DotWidth:         2,
					DotColorProvider: colourFunc,
				},
			},
			chart.AnnotationSeries{
				Annotations: annotations,
				Style: chart.Style{
					Show:             true,
					DotWidth:         1,
					DotColorProvider: colourFunc,
					TextLineSpacing:  5,
				},
			},
		},
	}

	fn := fmt.Sprintf("%s-%s.png", collectionID, time.Now().Truncate(time.Minute))
	f, err := os.Create(fn)
	if err != nil {
		return gerrors.Augment(err, "failed_to_create_scattergraph.create_file", map[string]string{
			"collection_id": collectionID,
			"vendor":        vendor.String(),
		})
	}
	defer f.Close()

	if err := graph.Render(chart.PNG, f); err != nil {
		slog.Error(ctx, "Failed to render scattergraph: %v", err)
		return gerrors.Augment(err, "failed_to_render_scattergraph", nil)
	}

	f, err = os.Open(fn)
	if err != nil {
		slog.Error(ctx, "Failed to read back file")
	}
	defer f.Close()

	if _, err := s.ChannelFileSend(m.ChannelID, fn, f); err != nil {
		slog.Error(ctx, "Failed to send scattergraph nft collection file to discord channel")
	}

	if err := os.Remove(fn); err != nil {
		slog.Error(ctx, "Failed to remove scattergraph of nft collection")
	}

	return nil
}

func convertVendor(vendor string) (solananftsproto.SolanaNFTVendor, error) {
	switch vendor {
	case "solanart":
		return solananftsproto.SolanaNFTVendor_SOLANART, nil
	case "magiceden":
		return solananftsproto.SolanaNFTVendor_MAGIC_EDEN, nil
	default:
		return solananftsproto.SolanaNFTVendor_UNKNOWN, gerrors.NotFound("solananft_vendor.not_found", nil)
	}
}

func convertCollectionID(c string) (string, error) {
	v, ok := collections[c]
	if !ok {
		return "", gerrors.NotFound("Colletion not found", map[string]string{
			"collection_id": c,
		})
	}

	return v, nil
}

func colourFunc(xr, yr chart.Range, index int, x, y float64) drawing.Color {
	return chart.Viridis(y, yr.GetMin(), yr.GetMax())
}
