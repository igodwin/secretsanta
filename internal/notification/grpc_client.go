package notification

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/igodwin/secretsanta/api/grpc/pb"
	"github.com/igodwin/secretsanta/pkg/participant"
)

// MessageTemplate defines the interface for generating notification messages
type MessageTemplate interface {
	Subject(giverName, recipientName string) string
	Body(giverName, recipientName string) string
}

// PapaElfTemplate is the default Secret Santa message template
type PapaElfTemplate struct{}

func (t *PapaElfTemplate) Subject(giverName, recipientName string) string {
	return "Your Official Secret Santa Assignment (Please Read Carefully)"
}

func (t *PapaElfTemplate) Body(giverName, recipientName string) string {
	return fmt.Sprintf(`Well, hello there %s,

After consulting the Official Elf Registry and cross-referencing it with the Nice List, twice, I might add, you have been selected to find a gift for %s this year.

Now, a few things you should know:

First, you're the only one who knows about this assignment, so... let's keep it that way, shall we? The whole "secret" part is rather crucial to the "Secret Santa" concept. I realize that sounds obvious, but you'd be surprised.

Second, don't delete this email right away. My memory isn't quite what it was back in, oh, 1823, and I won't remember who you're shopping for if you ask me later. So save it. Print it. Tattoo it on your arm if you have to.

Third, and this is important, please make a gift wish list available to the group. Your Secret Santa is counting on you. Think of it as helping them avoid the same panic I saw in Section C of the workshop last Tuesday. It wasn't pretty.

Now then, off you go. And remember: the best way to spread Christmas cheer is giving someone exactly what they asked for on their list.

With warm regards and slight concern for your organizational skills,

Papa Elf

North Pole Elf Personnel Director (Retired)
Secret Santa Coordinator (Current)`,
		giverName, recipientName)
}

type GRPCNotifier struct {
	client   pb.NotifierServiceClient
	conn     *grpc.ClientConn
	template MessageTemplate
	apiKey   string
}

func NewGRPCNotifier(serverAddr string) (*GRPCNotifier, error) {
	return NewGRPCNotifierWithAPIKey(serverAddr, "", &PapaElfTemplate{})
}

// NewGRPCNotifierWithAPIKey creates a notifier with an optional API key
func NewGRPCNotifierWithAPIKey(serverAddr string, apiKey string, template MessageTemplate) (*GRPCNotifier, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notifier service: %w", err)
	}

	client := pb.NewNotifierServiceClient(conn)

	if template == nil {
		template = &PapaElfTemplate{}
	}

	return &GRPCNotifier{
		client:   client,
		conn:     conn,
		template: template,
		apiKey:   apiKey,
	}, nil
}

// NewGRPCNotifierWithTemplate creates a notifier with a custom message template
func NewGRPCNotifierWithTemplate(serverAddr string, template MessageTemplate) (*GRPCNotifier, error) {
	return NewGRPCNotifierWithAPIKey(serverAddr, "", template)
}

func (g *GRPCNotifier) Close() error {
	return g.conn.Close()
}

// contextWithAPIKey adds the API key to the context metadata if present
func (g *GRPCNotifier) contextWithAPIKey(ctx context.Context) context.Context {
	if g.apiKey == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", g.apiKey))
}

// parseNotificationType parses a notification type string that may include an account
// Format: "type" or "type:account" (e.g., "email:notify")
// Returns: (type, account)
func parseNotificationType(notifType string) (string, string) {
	parts := strings.SplitN(notifType, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

func (g *GRPCNotifier) SendNotification(p *participant.Participant, archiveEmail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = g.contextWithAPIKey(ctx)

	// Parse notification type and account
	notifType, account := parseNotificationType(p.NotificationType)

	notificationType := pb.NotificationType_NOTIFICATION_TYPE_STDOUT
	switch notifType {
	case "email":
		notificationType = pb.NotificationType_NOTIFICATION_TYPE_EMAIL
	case "slack":
		notificationType = pb.NotificationType_NOTIFICATION_TYPE_SLACK
	case "ntfy":
		notificationType = pb.NotificationType_NOTIFICATION_TYPE_NTFY
	}

	subject := g.template.Subject(p.Name, p.Recipient.Name)
	body := g.template.Body(p.Name, p.Recipient.Name)

	// Build recipients list - support multiple contact methods
	recipients := make([]string, len(p.ContactInfo))
	copy(recipients, p.ContactInfo)

	// Add metadata
	metadata := map[string]string{
		"participant_name": p.Name,
		"recipient_name":   p.Recipient.Name,
		"event_type":       "secret_santa",
	}

	// Build BCC list with archive email if provided
	var bcc []string
	if archiveEmail != "" {
		bcc = []string{archiveEmail}
	}

	req := &pb.SendNotificationRequest{
		Type:       notificationType,
		Account:    account, // Set the account if specified
		Priority:   pb.Priority_PRIORITY_NORMAL,
		Subject:    subject,
		Body:       body,
		Recipients: recipients,
		Bcc:        bcc,
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

func (g *GRPCNotifier) SendBatchNotifications(participants []*participant.Participant, archiveEmail string, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	ctx = g.contextWithAPIKey(ctx)

	var requests []*pb.SendNotificationRequest
	for _, p := range participants {
		// Parse notification type and account
		notifType, account := parseNotificationType(p.NotificationType)

		notificationType := pb.NotificationType_NOTIFICATION_TYPE_STDOUT
		switch notifType {
		case "email":
			notificationType = pb.NotificationType_NOTIFICATION_TYPE_EMAIL
		case "slack":
			notificationType = pb.NotificationType_NOTIFICATION_TYPE_SLACK
		case "ntfy":
			notificationType = pb.NotificationType_NOTIFICATION_TYPE_NTFY
		}

		subject := g.template.Subject(p.Name, p.Recipient.Name)
		body := g.template.Body(p.Name, p.Recipient.Name)

		// Build recipients list - support multiple contact methods
		recipients := make([]string, len(p.ContactInfo))
		copy(recipients, p.ContactInfo)

		// Add metadata
		metadata := map[string]string{
			"participant_name": p.Name,
			"recipient_name":   p.Recipient.Name,
			"event_type":       "secret_santa",
		}

		// Build BCC list with archive email if provided
		var bcc []string
		if archiveEmail != "" {
			bcc = []string{archiveEmail}
		}

		req := &pb.SendNotificationRequest{
			Type:        notificationType,
			Account:     account, // Set the account if specified
			Priority:    pb.Priority_PRIORITY_NORMAL,
			Subject:     subject,
			Body:        body,
			Recipients:  recipients,
			Bcc:         bcc,
			Metadata:    metadata,
			ContentType: contentType,
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

// GetNotifiers queries the notifier service for available notification types
func (g *GRPCNotifier) GetNotifiers() ([]*pb.NotifierInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = g.contextWithAPIKey(ctx)

	resp, err := g.client.GetNotifiers(ctx, &pb.GetNotifiersRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get notifiers: %w", err)
	}

	return resp.Notifiers, nil
}
