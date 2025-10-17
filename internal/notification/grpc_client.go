package notification

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/igodwin/secretsanta/api/grpc/pb"
	"github.com/igodwin/secretsanta/pkg/participant"
)

type GRPCNotifier struct {
	client pb.NotifierServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCNotifier(serverAddr string) (*GRPCNotifier, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notifier service: %w", err)
	}

	client := pb.NewNotifierServiceClient(conn)
	
	return &GRPCNotifier{
		client: client,
		conn:   conn,
	}, nil
}

func (g *GRPCNotifier) Close() error {
	return g.conn.Close()
}

func (g *GRPCNotifier) SendNotification(p *participant.Participant, archiveEmail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	notificationType := pb.NotificationType_NOTIFICATION_TYPE_STDOUT
	switch p.NotificationType {
	case "email":
		notificationType = pb.NotificationType_NOTIFICATION_TYPE_EMAIL
	case "slack":
		notificationType = pb.NotificationType_NOTIFICATION_TYPE_SLACK
	}

	subject := "Secret Santa Assignment"
	body := fmt.Sprintf("Hi %s,\n\nYou have been assigned to give a gift to: %s\n\nHappy gifting!", 
		p.Name, p.Recipient.Name)

	// Build recipients list - support multiple contact methods
	recipients := make([]string, len(p.ContactInfo))
	copy(recipients, p.ContactInfo)

	// Add archive email as BCC if provided
	metadata := map[string]string{
		"participant_name": p.Name,
		"recipient_name":   p.Recipient.Name,
		"event_type":       "secret_santa",
	}
	if archiveEmail != "" {
		metadata["bcc"] = archiveEmail
	}

	req := &pb.SendNotificationRequest{
		Type:       notificationType,
		Priority:   pb.Priority_PRIORITY_NORMAL,
		Subject:    subject,
		Body:       body,
		Recipients: recipients,
		Metadata:   metadata,
	}

	resp, err := g.client.SendNotification(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	if !resp.Result.Success {
		return fmt.Errorf("notification failed: %s", resp.Result.Error)
	}

	log.Printf("Notification sent successfully to %v (ID: %s)", recipients, resp.Result.NotificationId)
	return nil
}

func (g *GRPCNotifier) SendBatchNotifications(participants []*participant.Participant, archiveEmail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var requests []*pb.SendNotificationRequest
	for _, p := range participants {
		notificationType := pb.NotificationType_NOTIFICATION_TYPE_STDOUT
		switch p.NotificationType {
		case "email":
			notificationType = pb.NotificationType_NOTIFICATION_TYPE_EMAIL
		case "slack":
			notificationType = pb.NotificationType_NOTIFICATION_TYPE_SLACK
		}

		subject := "Secret Santa Assignment"
		body := fmt.Sprintf("Hi %s,\n\nYou have been assigned to give a gift to: %s\n\nHappy gifting!", 
			p.Name, p.Recipient.Name)

		// Build recipients list - support multiple contact methods
		recipients := make([]string, len(p.ContactInfo))
		copy(recipients, p.ContactInfo)

		// Add archive email as BCC if provided
		metadata := map[string]string{
			"participant_name": p.Name,
			"recipient_name":   p.Recipient.Name,
			"event_type":       "secret_santa",
		}
		if archiveEmail != "" {
			metadata["bcc"] = archiveEmail
		}

		req := &pb.SendNotificationRequest{
			Type:       notificationType,
			Priority:   pb.Priority_PRIORITY_NORMAL,
			Subject:    subject,
			Body:       body,
			Recipients: recipients,
			Metadata:   metadata,
		}
		requests = append(requests, req)
	}

	batchReq := &pb.SendBatchNotificationsRequest{
		Notifications: requests,
	}

	resp, err := g.client.SendBatchNotifications(ctx, batchReq)
	if err != nil {
		return fmt.Errorf("failed to send batch notifications: %w", err)
	}

	for i, result := range resp.Results {
		if !result.Success {
			log.Printf("Failed to send notification to %v: %s", participants[i].ContactInfo, result.Error)
		} else {
			log.Printf("Notification sent successfully to %v (ID: %s)", participants[i].ContactInfo, result.NotificationId)
		}
	}

	return nil
}