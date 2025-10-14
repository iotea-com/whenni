package minio

import (
	"bytes"
	"context"
	"fmt"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	nodeHelpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/shared"
	"github.com/iotea-com/iotea/libs/engine/template"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (n *MinioActionSubnode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			if err := n.handle(ioData); err != nil {
				return node.Error{
					Type:   node.FatalError,
					Reason: err.Error(),
				}
			} else {
				n.notifyChannel <- node.Notification{
					Type:    node.NotifyDataProcessed,
					DataCtx: ioData.Ctx,
				}
			}
		}
	}

	return node.Error{
		Type: node.NoError,
	}
}

func (n *MinioActionSubnode) handle(ioData node.IoData) error {
	// Attempt to detect content type
	data, contentType, err := nodeHelpers.DetectDataContentType(ioData)
	if err != nil {
		return fmt.Errorf("unable to detect data content type: %v", err)
	}

	// Attempt to upload to Minio
	if err := n.writeDataToMinio(ioData.Ctx, data, contentType); err != nil {
		return fmt.Errorf("unable to upload data to MinIO: %v", err)

	}

	return nil
}

func (n *MinioActionSubnode) writeDataToMinio(dataCtx context.Context, data []byte, contentType string) error {
	// Initialize MinIO client
	client, err := minio.New(n.config.ThingMinioBucket.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(n.config.ThingMinioBucket.MinioAccessKeyId, n.config.ThingMinioBucket.MinioSecretAccessKey, ""),
		Secure: n.config.ThingMinioBucket.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	// Create & verify template
	keyTemplate := template.Template{
		Value:         n.config.Key,
		ModelInput:    &n.config.ModelInput,
		DefaultValues: n.config.DefaultValues,
	}

	// Modify the key as necessary
	// Preprocess the filter template
	keyTemplate.Preprocess()

	// Log success
	n.logger.Info().Ctx(dataCtx).Msgf("Uploading to MinIO bucket %s with key %s", n.config.ThingMinioBucket.BucketName, keyTemplate.Value)

	// Upload the file to MinIO bucket
	_, err = client.PutObject(dataCtx, n.config.ThingMinioBucket.BucketName, keyTemplate.Value, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload data to MinIO: %v", err)
	}

	// Log success
	n.logger.Info().Ctx(dataCtx).Msgf("Data successfully uploaded to MinIO bucket %s with key %s", n.config.ThingMinioBucket.BucketName, keyTemplate.Value)

	return nil
}
