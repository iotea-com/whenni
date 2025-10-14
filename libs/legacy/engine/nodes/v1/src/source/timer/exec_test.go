package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	timerSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/timer/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

// TestExec ensures that the source node successfully receives and processes messages from a CRON scheduler.
func TestExec(t *testing.T) {
	t.Run("successfully executes source node", func(t *testing.T) {
		n := New()

		// Create input parameters with a cancellation context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define the configuration for the Timer source node
		config := timerSourceNodeConfig.TimerSourceNodeConfig{
			CronExpression: "*/1 * * * * *",
		}

		// Create channels for errors and notifications
		errChan := make(chan error)
		notifyChan := make(chan node.Notification)

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Call Init with the context, config, lifetime, and observability
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		// Create channel to store published messages
		numMessages := 2

		// Verify that the subscriber node received the correct number of messages
		receivedMessagesCount := 0

	Loop:
		for {
			select {
			case <-ctx.Done():
				// Context canceled, exit the loop
				break Loop
			case <-n.GetOutputChans(ctx)[0].Channel:
				receivedMessagesCount++

				// Verify that the subscriber node received the correct number of messages
				if receivedMessagesCount == numMessages {
					t.Logf("Message count sent(%v) and received(%v) MATCH", receivedMessagesCount, numMessages)
					break Loop
				}
			case notification := <-notifyChan:
				t.Logf("received notification %v", notification)
			case err := <-errChan:
				t.Fatalf("Exec() returned error: %v", err)
			}
		}

		if receivedMessagesCount != numMessages {
			t.Errorf("Number of received messages (%v) does not match expected number of messages (%v)", receivedMessagesCount, numMessages)
		} else {
			t.Logf("Message count sent(%v) and received(%v) MATCH", receivedMessagesCount, numMessages)
		}

		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node: %v", deinitErr)
		}
	})
}
