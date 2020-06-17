package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestS3Handler(t *testing.T) {
	event := "{\"Records\":[{\"eventVersion\":\"2.1\",\"eventSource\":\"aws:s3\",\"awsRegion\":\"us-west-2\",\"eventTime\":\"2020-06-17T17:18:02.044Z\",\"eventName\":\"ObjectCreated:Put\",\"userIdentity\":{\"principalId\":\"AMB55Y5FSTZ6B\"},\"requestParameters\":{\"sourceIPAddress\":\"67.161.57.201\"},\"responseElements\":{\"x-amz-id-2\":\"E9e82Tw/2XgXVHcFunipWv67t10N9UYRPqPRou7TQ3wYc7R+TClR0l1gwKTL9jclrg48Unpl2hm1ZK3a59164gAwaUM/0Pfe\",\"x-amz-request-id\":\"86084DA99A863547\"},\"s3\":{\"s3SchemaVersion\":\"1.0\",\"configurationId\":\"timeline-dev-timeline_s3_handler-68a737133cefb25bff959852b8f04754\",\"bucket\":{\"name\":\"timeline-moments\",\"ownerIdentity\":{\"principalId\":\"AMB55Y5FSTZ6B\"},\"arn\":\"arn:aws:s3:::timeline-moments\"},\"object\":{\"key\":\"20180403_133305.jpg\",\"size\":1247691,\"urlDecodedKey\":\"\",\"versionId\":\"\",\"eTag\":\"bfc1ed8cc49e2293e4ad2f08fbd633a9\",\"sequencer\":\"005EEA504E702958A9\"}}}]}"
	s3Event := new(events.S3Event)
	json.Unmarshal([]byte(event), s3Event)
	type args struct {
		ctx   context.Context
		event events.S3Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "temp",
			args:    args{event: *s3Event},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := S3Handler(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("S3Handler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestLocalJpg(t *testing.T) {
// 	f, err := os.Open("/home/hangc/ws/go/src/timeline_service/20180403_133305.jpg")
// 	if err != nil {
// 		log.Infof("err=%v", err)
// 	}
// 	defer f.Close()
// 	reader := bufio.NewReader(f)
// 	rawExif, err := exif_v2.SearchAndExtractExifWithReader(reader)
// 	im := exif.NewIfdMappingWithStandard()
// 	ti := exif.NewTagIndex()
// 	entries := make([]IfdEntry, 0)
// 	visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {

// 		ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)

// 		it, err := ti.Get(ifdPath, tagId)
// 		if err != nil {
// 			log.Errorf("err = %v", err)
// 		}

// 		valueString := ""
// 		var value interface{}
// 		if tagType.Type() == exif.TypeUndefined {
// 			var err error
// 			value, err = valueContext.Undefined()
// 			if err != nil {
// 				if err == exif.ErrUnhandledUnknownTypedTag {
// 					value = nil
// 				} else {
// 					log.Panic(err)
// 				}
// 			}

// 			valueString = fmt.Sprintf("%v", value)
// 		} else {
// 			valueString, err = valueContext.FormatFirst()
// 			log.Errorf("err = %v", err)

// 			value = valueString
// 		}

// 		entry := IfdEntry{
// 			IfdPath:     ifdPath,
// 			FqIfdPath:   fqIfdPath,
// 			IfdIndex:    ifdIndex,
// 			TagId:       tagId,
// 			TagName:     it.Name,
// 			TagTypeId:   tagType.Type(),
// 			TagTypeName: tagType.Name(),
// 			UnitCount:   valueContext.UnitCount(),
// 			Value:       value,
// 			ValueString: valueString,
// 		}

// 		entries = append(entries, entry)

// 		return nil
// 	}
// 	_, err = exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)
// 	for _, entry := range entries {
// 		log.Infof("tagName: %v, ValueString: %v, value = %v", entry.TagName, entry.ValueString, entry.Value)
// 	}
// 	getGpsInfo(rawExif)
// }

// func getGpsInfo(rawExif []byte) {
// 	im := exif_v2.NewIfdMapping()

// 	err := exif_v2.LoadStandardIfds(im)

// 	ti := exif_v2.NewTagIndex()

// 	_, index, err := exif_v2.Collect(im, ti, rawExif)

// 	ifd, err := index.RootIfd.ChildWithIfdPath(exifcommon.IfdGpsInfoStandardIfdIdentity)

// 	gi, err := ifd.GpsInfo()

// 	if err != nil {
// 		log.Errorf("error = %s", err)
// 	}

// 	log.Infof("%s\n", gi)
// }
