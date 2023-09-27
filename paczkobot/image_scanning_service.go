package paczkobot

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strings"

	"github.com/alufers/paczkobot/tghelpers"
	"github.com/fogleman/gg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/aztec"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type ScannedShipmentNumberCandidate struct {
	Number string
	Image  image.Image
}

type ImageScanningService struct {
	App *BotApp
}

func NewImageScanningService(app *BotApp) *ImageScanningService {
	return &ImageScanningService{
		App: app,
	}
}

func (i *ImageScanningService) OnUpdate(ctx context.Context) bool {
	update := tghelpers.UpdateFromCtx(ctx)
	if update.Message != nil && update.Message.Photo != nil && len(update.Message.Photo) > 0 {

		photo := update.Message.Photo[len(update.Message.Photo)-1]

		file, err := i.App.Bot.GetFile(tgbotapi.FileConfig{
			FileID: photo.FileID,
		})
		if err != nil {
			log.Printf("Failed to get file: %v", err)
			return false
		}
		url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", i.App.Bot.Token, file.FilePath)
		err = i.ScanIncomingImage(ctx, tghelpers.ArgsFromCtx(ctx), url)
		if err != nil {
			log.Printf("Failed to ScanIncomingImage: %v", err)
		}
		return true
	}
	return false
}

func (i *ImageScanningService) ScanIncomingImage(ctx context.Context, args *tghelpers.CommandArguments, url string) error {
	progress, err := tghelpers.NewProgressMessage(i.App.Bot, args.ChatID, "(1/2) Fetching image...")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	log.Printf("Fetched image from %s", url)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch image: %s", resp.Status)
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	img, err := jpeg.Decode(bytes.NewBuffer(buf.Bytes()))
	imgType := "jpeg"
	if err != nil {
		return fmt.Errorf("failed to decode image: %s", err)
	}

	err = progress.UpdateText(fmt.Sprintf("(2/3) Scanning %v image for barcodes...", imgType))
	if err != nil {
		return err
	}

	// barcode scanning
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return fmt.Errorf("failed to create binary bitmap from image: %s", err)
	}

	shipmentNumberCandidates := make([]ScannedShipmentNumberCandidate, 0)

	readerConstructors := []func() gozxing.Reader{
		oned.NewEAN13Reader,
		oned.NewEAN8Reader,
		oned.NewUPCAReader,
		oned.NewUPCEReader,
		oned.NewCode128Reader,
		oned.NewCode39Reader,
		oned.NewITFReader,
		oned.NewCodaBarReader,
		oned.NewCode93Reader,
		qrcode.NewQRCodeReader,
		func() gozxing.Reader { return aztec.NewAztecReader() },
		func() gozxing.Reader { return datamatrix.NewDataMatrixReader() },
	}
	for _, readerConstructor := range readerConstructors {
		reader := readerConstructor()
		readerName := fmt.Sprintf("%T", reader)
		readerName = strings.TrimPrefix(strings.TrimSuffix(readerName, "Reader"), "*oned.")
		err := progress.UpdateText(fmt.Sprintf("(2/2) Scanning %v image for barcodes... (using %v)", imgType, readerName))
		if err != nil {
			return fmt.Errorf("failed to update progress message: %s", err)
		}
		result, err := reader.Decode(bmp, map[gozxing.DecodeHintType]interface{}{
			gozxing.DecodeHintType_TRY_HARDER: true,
		})
		if err == nil {
			img2 := i.DrawResultPoints(img, result.GetResultPoints())

			shipmentNumberCandidates = append(shipmentNumberCandidates, ScannedShipmentNumberCandidate{
				Number: result.GetText(),
				Image:  img2,
			})
			// return nil
		}
	}

	if err := progress.Delete(); err != nil {
		log.Printf("failed to delete progress message: %s", err)
	}

	for _, candidate := range shipmentNumberCandidates {

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, candidate.Image, nil)
		if err != nil {
			return fmt.Errorf("failed to encode image: %s", err)
		}
		msg := tgbotapi.NewPhoto(args.ChatID, tgbotapi.FileBytes{
			Name:  "result.jpg",
			Bytes: buf.Bytes(),
		})
		msg.Caption = fmt.Sprintf("Found barcode: %v", candidate.Number)
		_, err = i.App.Bot.Send(msg)
		if err != nil {
			return err
		}
		args.NamedArguments = map[string]string{
			"shipmentNumber": candidate.Number,
		}

		err := (&TrackCommand{App: i.App}).Execute(context.WithValue(
			ctx,
			tghelpers.ArgsContextKey,
			args,
		))
		if err != nil {
			return err
		}
	}

	return nil
}

func (*ImageScanningService) DrawResultPoints(img image.Image, points []gozxing.ResultPoint) image.Image {
	if len(points) <= 1 {
		return img
	}
	ctx := gg.NewContextForImage(img)
	ctx.SetRGB(1, 0, 0)
	minX, minY, maxX, maxY := ctx.Width(), ctx.Height(), 0, 0
	ctx.MoveTo(points[0].GetX(), points[0].GetY())
	for _, point := range points {
		x, y := int(point.GetX()), int(point.GetY())
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		ctx.LineTo(float64(x), float64(y))
	}
	if len(points) > 2 {
		ctx.LineTo(points[0].GetX(), points[0].GetY())
	}
	ctx.SetLineWidth(3)
	ctx.Stroke()
	// crop the image to the bounding box with a margin
	margin := int(float64(max(maxX-minX, maxY-minY)) * 0.2)
	minX = max(0, minX-margin)
	minY = max(0, minY-margin)
	maxX = min(ctx.Width(), maxX+margin)
	maxY = min(ctx.Height(), maxY+margin)
	ctx2 := gg.NewContext(maxX-minX, maxY-minY)
	ctx2.DrawImage(ctx.Image(), -minX, -minY)

	return ctx2.Image()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
