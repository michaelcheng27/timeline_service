package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"timeline_service/src/cmd"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dsoprea/go-exif"
	exif_v2 "github.com/dsoprea/go-exif/v2"
	exifcommon "github.com/dsoprea/go-exif/v2/common"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func S3Handler(ctx context.Context, event events.S3Event) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// 3) Create a new AWS S3 downloader
	downloader := s3manager.NewDownloader(sess)
	log.Infof("event = %v", event)
	b, _ := json.Marshal(event)
	log.Infof("event = %v", string(b))
	for _, s3Record := range event.Records {
		err := processS3Event(s3Record, downloader)
		panicIf(err)
	}
	return nil
}

func processS3Event(record events.S3EventRecord, downloader *s3manager.Downloader) error {
	buff := &aws.WriteAtBuffer{}
	s3Identity := record.S3
	bucket, object := s3Identity.Bucket, s3Identity.Object
	key := object.Key
	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket.Name),
			Key:    aws.String(key),
		})
	if err != nil {
		log.Errorf("Unable to download bucket %v, key: %v, %v", bucket.Name, key, err)
		return err
	}
	log.Infof("Downloaded: bucket = %v, key = %v, bytes %v", bucket.Name, key, numBytes)
	reader := bytes.NewReader(buff.Bytes())
	rawExif, err := exif_v2.SearchAndExtractExifWithReader(reader)
	panicIf(err)
	moment, err := getMoment(rawExif)
	panicIf(err)
	moment.S3Bucket = bucket.Name
	moment.S3Key = key

	log.Infof("moment = %v", moment)
	err = cmd.PutMoment(moment)
	return err
}

func getGpsInfo(rawExif []byte) (*exif_v2.GpsInfo, error) {
	im := exif_v2.NewIfdMapping()

	err := exif_v2.LoadStandardIfds(im)
	panicIf(err)

	ti := exif_v2.NewTagIndex()

	_, index, err := exif_v2.Collect(im, ti, rawExif)
	panicIf(err)

	ifd, err := index.RootIfd.ChildWithIfdPath(exifcommon.IfdGpsInfoStandardIfdIdentity)
	panicIf(err)

	return ifd.GpsInfo()
}

type IfdEntry struct {
	IfdPath     string                `json:"ifd_path"`
	FqIfdPath   string                `json:"fq_ifd_path"`
	IfdIndex    int                   `json:"ifd_index"`
	TagId       uint16                `json:"tag_id"`
	TagName     string                `json:"tag_name"`
	TagTypeId   exif.TagTypePrimitive `json:"tag_type_id"`
	TagTypeName string                `json:"tag_type_name"`
	UnitCount   uint32                `json:"unit_count"`
	Value       interface{}           `json:"value"`
	ValueString string                `json:"value_string"`
}

func getMoment(rawExif []byte) (cmd.Moment, error) {
	gpsInfo, err := getGpsInfo(rawExif)
	moment := cmd.Moment{
		Latitude:  gpsInfo.Latitude.Decimal(),
		Longitude: gpsInfo.Longitude.Decimal(),
		Id:        uuid.New().String(),
	}

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()
	entries := make([]IfdEntry, 0)
	visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {

		ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)
		panicIf(err)
		it, err := ti.Get(ifdPath, tagId)
		panicIf(err)

		valueString := ""
		var value interface{}
		if tagType.Type() == exif.TypeUndefined {
			var err error
			value, err = valueContext.Undefined()
			if err != nil {
				if err == exif.ErrUnhandledUnknownTypedTag {
					value = nil
				} else {
					panicIf(err)
				}
			}

			valueString = fmt.Sprintf("%v", value)
		} else {
			valueString, err = valueContext.FormatFirst()
			panicIf(err)

			value = valueString
		}
		log.Infof("tagName = %v", it.Name)
		switch tagName := it.Name; tagName {
		case "ImageUniqueID":
			moment.ImageId = valueString
		case "Make":
			moment.CameraMake = valueString
		case "Model":
			moment.CameraModel = valueString
		case "ImageWidth":
			moment.Width = valueString
		case "ImageLength":
			moment.Height = valueString
		case "DateTime":
			moment.TimeTaken = valueString
		}

		entry := IfdEntry{
			IfdPath:     ifdPath,
			FqIfdPath:   fqIfdPath,
			IfdIndex:    ifdIndex,
			TagId:       tagId,
			TagName:     it.Name,
			TagTypeId:   tagType.Type(),
			TagTypeName: tagType.Name(),
			UnitCount:   valueContext.UnitCount(),
			Value:       value,
			ValueString: valueString,
		}
		log.Infof("entry : %v", entry)

		entries = append(entries, entry)

		return nil
	}
	_, err = exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)
	return moment, err
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Info("Hello S3 handler")
	lambda.Start(S3Handler)
}
